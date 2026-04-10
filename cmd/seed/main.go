package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

const (
	adminUUID = "11111111-1111-1111-1111-111111111111"
	userUUID  = "22222222-2222-2222-2222-222222222222"

	slotDuration = 30 * time.Minute
	seedDays     = 7 // количество дней вперёд для генерации слотов
)

func main() {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		dsn = "postgres://postgres:postgres@localhost:5432/booking?sslmode=disable" //nolint:gosec
	}

	ctx := context.Background()

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatalf("не удалось подключиться к БД: %v", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		log.Fatalf("БД недоступна: %v", err)
	}
	defer pool.Close()

	seedUsers(ctx, pool)
	rooms := seedRooms(ctx, pool)
	schedules := seedSchedules(ctx, pool, rooms)
	seedSlots(ctx, pool, schedules)

	log.Println("seed завершён успешно")
}

// Пользователи

type seedUser struct {
	id    string
	email string
	pass  string
	role  string
}

func seedUsers(ctx context.Context, pool *pgxpool.Pool) {
	users := []seedUser{
		{adminUUID, "admin@example.com", "admin123", "admin"},
		{userUUID, "user@example.com", "user123", "user"},
	}

	for _, u := range users {
		hash, err := bcrypt.GenerateFromPassword([]byte(u.pass), bcrypt.DefaultCost)
		if err != nil {
			log.Fatalf("bcrypt для %s: %v", u.email, err)
		}

		_, err = pool.Exec(ctx,
			`INSERT INTO users (id, email, password, role)
			 VALUES ($1, $2, $3, $4)
			 ON CONFLICT (id) DO NOTHING`,
			u.id, u.email, string(hash), u.role,
		)
		if err != nil {
			log.Fatalf("вставка пользователя %s: %v", u.email, err)
		}
		log.Printf("пользователь: %s (%s)", u.email, u.role)
	}
}

// Переговорки

type roomSeed struct {
	id          uuid.UUID
	name        string
	description string
	capacity    int
}

func seedRooms(ctx context.Context, pool *pgxpool.Pool) []roomSeed {
	rooms := []roomSeed{
		{uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"), "Переговорка «Альфа»", "Большая переговорка на 1 этаже", 10},
		{uuid.MustParse("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"), "Переговорка «Бета»", "Маленькая переговорка на 2 этаже", 4},
		{uuid.MustParse("cccccccc-cccc-cccc-cccc-cccccccccccc"), "Переговорка «Гамма»", "Переговорка без окон", 6},
	}

	for _, rm := range rooms {
		_, err := pool.Exec(ctx,
			`INSERT INTO rooms (id, name, description, capacity)
			 VALUES ($1, $2, $3, $4)
			 ON CONFLICT (id) DO NOTHING`,
			rm.id, rm.name, rm.description, rm.capacity,
		)
		if err != nil {
			log.Fatalf("вставка переговорки %s: %v", rm.name, err)
		}
		log.Printf("переговорка: %s", rm.name)
	}

	return rooms
}

// Расписания

type scheduleSeed struct {
	id         uuid.UUID
	roomID     uuid.UUID
	daysOfWeek []int
	startTime  time.Time
	endTime    time.Time
}

func seedSchedules(ctx context.Context, pool *pgxpool.Pool, rooms []roomSeed) []scheduleSeed {
	weekdays := []int{1, 2, 3, 4, 5}
	schedules := []scheduleSeed{
		{
			id:         uuid.MustParse("dddddddd-dddd-dddd-dddd-dddddddddddd"),
			roomID:     rooms[0].id,
			daysOfWeek: weekdays,
			startTime:  timeOnly(9),
			endTime:    timeOnly(18),
		},
		{
			id:         uuid.MustParse("eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee"),
			roomID:     rooms[1].id,
			daysOfWeek: weekdays,
			startTime:  timeOnly(10),
			endTime:    timeOnly(16),
		},
		{
			id:         uuid.MustParse("ffffffff-ffff-ffff-ffff-ffffffffffff"),
			roomID:     rooms[2].id,
			daysOfWeek: []int{1, 2, 3, 4, 5, 6}, // пн-сб
			startTime:  timeOnly(8),
			endTime:    timeOnly(20),
		},
	}

	for _, s := range schedules {
		_, err := pool.Exec(ctx,
			`INSERT INTO schedules (id, room_id, start_time, end_time)
			 VALUES ($1, $2, $3, $4)
			 ON CONFLICT (id) DO NOTHING`,
			s.id, s.roomID,
			s.startTime.Format("15:04:05"),
			s.endTime.Format("15:04:05"),
		)
		if err != nil {
			log.Fatalf("вставка расписания для %s: %v", s.roomID, err)
		}

		for _, day := range s.daysOfWeek {
			_, err := pool.Exec(ctx,
				`INSERT INTO schedule_days (schedule_id, day_of_week)
				 VALUES ($1, $2)
				 ON CONFLICT DO NOTHING`,
				s.id, day,
			)
			if err != nil {
				log.Fatalf("вставка дня расписания %d: %v", day, err)
			}
		}
		log.Printf("расписание: комната %s, дни %v, %s–%s",
			s.roomID,
			s.daysOfWeek,
			s.startTime.Format("15:04"),
			s.endTime.Format("15:04"),
		)
	}

	return schedules
}

// --- Слоты ---

func seedSlots(ctx context.Context, pool *pgxpool.Pool, schedules []scheduleSeed) {
	today := time.Now().UTC().Truncate(24 * time.Hour)
	count := 0

	for _, s := range schedules {
		for d := range seedDays {
			day := today.AddDate(0, 0, d)
			isoWeekday := int(day.Weekday())
			if isoWeekday == 0 {
				isoWeekday = 7
			}

			if !containsDay(s.daysOfWeek, isoWeekday) {
				continue
			}

			slotStart := time.Date(day.Year(), day.Month(), day.Day(),
				s.startTime.Hour(), s.startTime.Minute(), 0, 0, time.UTC)
			slotEnd := time.Date(day.Year(), day.Month(), day.Day(),
				s.endTime.Hour(), s.endTime.Minute(), 0, 0, time.UTC)

			for cur := slotStart; cur.Add(slotDuration).Before(slotEnd) || cur.Add(slotDuration).Equal(slotEnd); cur = cur.Add(slotDuration) {
				_, err := pool.Exec(ctx,
					`INSERT INTO slots (id, room_id, start_time, end_time)
					 VALUES ($1, $2, $3, $4)
					 ON CONFLICT DO NOTHING`,
					uuid.New(), s.roomID, cur, cur.Add(slotDuration),
				)
				if err != nil {
					log.Fatalf("вставка слота: %v", err)
				}
				count++
			}
		}
	}

	log.Printf("слоты: сгенерировано %d слотов на %d дней вперёд", count, seedDays)
}

func timeOnly(hour int) time.Time {
	return time.Date(0, 1, 1, hour, 0, 0, 0, time.UTC)
}

func containsDay(days []int, day int) bool {
	for _, d := range days {
		if d == day {
			return true
		}
	}
	return false
}

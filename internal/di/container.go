package di

import "context"

type Container struct {
	db       dbProvider
	configs  configProvider
	repos    repositoryProvider
	services serviceProvider
	handlers handlerProvider
	clients  clientProvider
}

func NewContainer() (*Container, error) {
	ctx := context.Background()

	c := &Container{}

	c.init(ctx)

	return c, nil
}

func (c *Container) init(ctx context.Context) {
	inits := []func(context.Context){
		c.initConfigs,
		c.initDB,
		c.initRepos,
		c.initServices,
		c.initHandlers,
	}

	for _, fn := range inits {
		fn(ctx)
	}
}

func (c *Container) initConfigs(_ context.Context) {
	c.JWTConfig()
	c.DBConfig()
}

func (c *Container) initDB(ctx context.Context) {
	c.DB(ctx)
}

func (c *Container) initRepos(ctx context.Context) {
	c.BookingRepo(ctx)
	c.RoomRepo(ctx)
	c.ScheduleRepo(ctx)
	c.SlotRepo(ctx)
	c.UserRepo(ctx)
}

func (c *Container) initServices(ctx context.Context) {
	c.AuthService(ctx)
	c.BookingService(ctx)
	c.RoomService(ctx)
	c.ScheduleService(ctx)
	c.SlotService(ctx)
}

func (c *Container) initHandlers(ctx context.Context) {
	c.AuthHandler(ctx)
	c.BookingHandler(ctx)
	c.RoomHandler(ctx)
	c.ScheduleHandler(ctx)
	c.SlotHandler(ctx)
}

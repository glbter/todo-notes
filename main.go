package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
	"todoNote/internal/repo/postgres"
	http2 "todoNote/internal/server/http"
)

const(
	httpPortEnv = "HTTP_PORT"
)

func testHandler(w http.ResponseWriter, r *http.Request) {
	time.Sleep(time.Millisecond*100)
	fmt.Fprintln(w,"hello there!")
}

func main() {
	//mux := http.NewServeMux()
	//mux.HandleFunc("/", testHandler)
	//log.Fatal(http.ListenAndServe(":8080", apmhttp.Wrap(mux)))
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		call:= <-c
		log.Printf("received signal to shutdown: %v", call)
		cancel()
	} ()

	dbConnection := make(chan *pgx.Conn, 1)
	go func(ctx context.Context,c chan *pgx.Conn) {
		defer close(c)
		conn, err := postgres.Connect(ctx)
		if err != nil {
			log.Fatal(err)
		}
		c <- conn
	} (ctx, dbConnection)

	conn := <- dbConnection
	defer conn.Close(context.Background())

	go func(ctx context.Context, conn *pgx.Conn) {
		if err := postgres.Ping(ctx, conn); err != nil {
			log.Fatal(err)
		}
	}(ctx, conn)

	log.Println("successfully connected to db")
	repos := http2.Repositories{
		Note: postgres.NewRepoNote(conn),
		User: postgres.NewRepoUser(conn),
	}

	r, err := http2.NewRouter(repos)
	if err != nil {
		log.Fatal(err)
	}

	srv := &http.Server{
		Handler: r,
		Addr: ":"+os.Getenv(httpPortEnv),
	}

	go func() {
		log.Println("starting http server")
		log.Fatal(srv.ListenAndServe())
	}()

	<-ctx.Done()

	ctxShutdown, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctxShutdown); err != nil {
		log.Fatalf("could not shutdown properly: %v", err)
	}

	log.Printf("server shut down")
}

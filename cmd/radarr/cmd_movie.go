package main

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/SkYNewZ/radarr"
	"github.com/urfave/cli/v2"
)

var movieCommand *cli.Command = &cli.Command{
	Name:        "movies",
	Usage:       "Perform actions on movies",
	Description: "List movies, get a single movie",
}

func init() {
	app.Commands = append(app.Commands, movieCommand)
	movieCommand.Subcommands = append(movieCommand.Subcommands, []*cli.Command{
		{
			Name:    "list",
			Usage:   "List all movies in your collection",
			Aliases: []string{"ls"},
			Action:  listMovies,
		},
		{
			Name:   "get",
			Usage:  "Search a movie by ID",
			Action: getMovie,
		},
		{
			Name:   "delete",
			Usage:  "Delete a movie by ID",
			Action: deleteMovie,
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:  "delete-files",
					Usage: "If true the movie folder and all files will be deleted when the movie is deleted",
				},
				&cli.BoolFlag{
					Name:  "blacklist",
					Usage: "If true the movie TMDB ID will be added to the import exclusions list when the movie is deleted",
				},
			},
		},
		{
			Name:   "upcoming",
			Usage:  "List upcoming movies",
			Action: upcoming,
			Flags: []cli.Flag{
				&cli.TimestampFlag{
					Name:   "start",
					Usage:  "Specify a start date",
					Layout: "2006-01-02T15:04:05Z",
				},
				&cli.TimestampFlag{
					Name:   "end",
					Usage:  "Specify a end date",
					Layout: "2006-01-02T15:04:05Z",
				},
			},
		},
		{
			Name:   "excluded",
			Usage:  "Gets movies marked as List Exclusions",
			Action: excluded,
		},
	}...)
}

func listMovies(*cli.Context) error {
	movies, err := radarrClient.Movies.List()
	if err != nil {
		return err
	}

	t.SetHeader([]string{"Id", "Title", "Downloaded", "Monitored", "Added"})
	var data [][]string

	// Filter movies
	for _, movie := range movies {
		data = append(data, []string{
			strconv.Itoa(movie.ID),
			movie.Title,
			strconv.FormatBool(movie.Downloaded),
			strconv.FormatBool(movie.Monitored),
			movie.Added.Format(time.RFC3339),
		})
	}

	t.AppendBulk(data)
	t.Render()
	return nil
}

func getMovie(c *cli.Context) error {
	movieID := c.Args().First()
	if movieID == "" {
		return errors.New("Please specify a movie ID")
	}

	m, err := strconv.Atoi(movieID)
	if err != nil {
		return err
	}

	movie, err := radarrClient.Movies.Get(m)
	if err != nil {
		return err
	}

	v := reflect.ValueOf(*movie)
	typeOfMovie := v.Type()

	var data [][]string
	for i := 0; i < v.NumField(); i++ {
		switch typeOfMovie.Field(i).Name {
		case "Overview":
		case "Website":
		case "SortTitle":
		case "FolderName":
		case "CleanTitle":
		case "TitleSlug":
		default:
			if value, ok := v.Field(i).Interface().(string); ok {
				data = append(data, []string{
					fmt.Sprintf("%s:", typeOfMovie.Field(i).Name),
					value,
				})
			}
		}
	}

	t.SetCenterSeparator("")
	t.SetColumnSeparator("")
	t.SetRowSeparator("")
	t.AppendBulk(data)
	t.Render()
	return nil
}

func upcoming(c *cli.Context) error {
	start := c.Value("start").(cli.Timestamp)
	end := c.Value("end").(cli.Timestamp)

	opts := &radarr.UpcomingOptions{}
	if start.Value() != nil {
		opts.Start = start.Value()
	}
	if end.Value() != nil {
		opts.End = end.Value()
	}

	movies, err := radarrClient.Movies.Upcoming(opts)
	if err != nil {
		return err
	}

	t.SetHeader([]string{"Id", "Title", "Downloaded", "Monitored", "Added"})
	var data [][]string

	// Filter movies
	for _, movie := range movies {
		data = append(data, []string{
			strconv.Itoa(movie.ID),
			movie.Title,
			strconv.FormatBool(movie.Downloaded),
			strconv.FormatBool(movie.Monitored),
			movie.Added.Format(time.RFC3339),
		})
	}

	t.AppendBulk(data)
	t.Render()
	return nil
}

func deleteMovie(c *cli.Context) error {
	movieID := c.Args().First()
	if movieID == "" {
		return errors.New("Please specify a movie ID")
	}

	m, err := strconv.Atoi(movieID)
	if err != nil {
		return err
	}

	radarrMovie, err := radarrClient.Movies.Get(m)
	if err != nil {
		return err
	}

	err = radarrClient.Movies.Delete(radarrMovie, &radarr.DeleteMovieOptions{
		AddExclusion: c.Bool("blacklist"),
		DeleteFiles:  c.Bool("delete-files"),
	})
	if err != nil {
		return err
	}

	fmt.Println("Successfully deleted")
	return nil
}

func excluded(c *cli.Context) error {
	movies, err := radarrClient.Movies.Excluded()
	if err != nil {
		return err
	}

	// Set header
	var headers []string
	v := reflect.ValueOf(*movies[0])
	typeOfMovie := v.Type()
	for i := 0; i < v.NumField(); i++ {
		headers = append(headers, typeOfMovie.Field(i).Name)
	}

	var data [][]string
	for _, movie := range movies {
		v := reflect.ValueOf(*movie)

		var row []string
		for i := 0; i < v.NumField(); i++ {
			row = append(row, fmt.Sprintf("%v", v.Field(i).Interface()))
		}
		data = append(data, row)
	}

	t.SetHeader(headers)
	t.AppendBulk(data)
	t.Render()
	return nil
}

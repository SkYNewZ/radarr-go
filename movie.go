package radarr

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// Movie Radarr movie
type Movie struct {
	Title             string `json:"title"`
	AlternativeTitles []struct {
		SourceType string `json:"sourceType"`
		MovieID    int    `json:"movieId"`
		Title      string `json:"title"`
		SourceID   int    `json:"sourceId"`
		Votes      int    `json:"votes"`
		VoteCount  int    `json:"voteCount"`
		Language   string `json:"language"`
		ID         int    `json:"id"`
	} `json:"alternativeTitles"`
	SecondaryYearSourceID int       `json:"secondaryYearSourceId"`
	SortTitle             string    `json:"sortTitle"`
	SizeOnDisk            int64     `json:"sizeOnDisk"`
	Status                string    `json:"status"`
	Overview              string    `json:"overview"`
	InCinemas             time.Time `json:"inCinemas"`
	PhysicalRelease       time.Time `json:"physicalRelease"`
	Images                []struct {
		CoverType string `json:"coverType"`
		URL       string `json:"url"`
	} `json:"images"`
	Website             string    `json:"website"`
	Downloaded          bool      `json:"downloaded"`
	Year                int       `json:"year"`
	HasFile             bool      `json:"hasFile"`
	YouTubeTrailerID    string    `json:"youTubeTrailerId"`
	Studio              string    `json:"studio"`
	Path                string    `json:"path"`
	ProfileID           int       `json:"profileId"`
	Monitored           bool      `json:"monitored"`
	MinimumAvailability string    `json:"minimumAvailability"`
	IsAvailable         bool      `json:"isAvailable"`
	FolderName          string    `json:"folderName"`
	Runtime             int       `json:"runtime"`
	LastInfoSync        time.Time `json:"lastInfoSync"`
	CleanTitle          string    `json:"cleanTitle"`
	ImdbID              string    `json:"imdbId"`
	TmdbID              int       `json:"tmdbId"`
	TitleSlug           string    `json:"titleSlug"`
	Genres              []string  `json:"genres"`
	Tags                []int     `json:"tags"`
	Added               time.Time `json:"added"`
	Ratings             struct {
		Votes int     `json:"votes"`
		Value float64 `json:"value"`
	} `json:"ratings"`
	MovieFile struct {
		MovieID      int       `json:"movieId"`
		RelativePath string    `json:"relativePath"`
		Size         int64     `json:"size"`
		DateAdded    time.Time `json:"dateAdded"`
		SceneName    string    `json:"sceneName"`
		Quality      struct {
			Quality  Quality `json:"quality"`
			Revision struct {
				Version  int  `json:"version"`
				Real     int  `json:"real"`
				IsRepack bool `json:"isRepack"`
			} `json:"revision"`
		} `json:"quality"`
		Edition   string `json:"edition"`
		MediaInfo struct {
			ContainerFormat              string  `json:"containerFormat"`
			VideoFormat                  string  `json:"videoFormat"`
			VideoCodecID                 string  `json:"videoCodecID"`
			VideoProfile                 string  `json:"videoProfile"`
			VideoCodecLibrary            string  `json:"videoCodecLibrary"`
			VideoBitrate                 int     `json:"videoBitrate"`
			VideoBitDepth                int     `json:"videoBitDepth"`
			VideoMultiViewCount          int     `json:"videoMultiViewCount"`
			VideoColourPrimaries         string  `json:"videoColourPrimaries"`
			VideoTransferCharacteristics string  `json:"videoTransferCharacteristics"`
			Width                        int     `json:"width"`
			Height                       int     `json:"height"`
			AudioFormat                  string  `json:"audioFormat"`
			AudioCodecID                 string  `json:"audioCodecID"`
			AudioCodecLibrary            string  `json:"audioCodecLibrary"`
			AudioAdditionalFeatures      string  `json:"audioAdditionalFeatures"`
			AudioBitrate                 int     `json:"audioBitrate"`
			RunTime                      string  `json:"runTime"`
			AudioStreamCount             int     `json:"audioStreamCount"`
			AudioChannels                int     `json:"audioChannels"`
			AudioChannelPositions        string  `json:"audioChannelPositions"`
			AudioChannelPositionsText    string  `json:"audioChannelPositionsText"`
			AudioProfile                 string  `json:"audioProfile"`
			VideoFps                     float64 `json:"videoFps"`
			AudioLanguages               string  `json:"audioLanguages"`
			Subtitles                    string  `json:"subtitles"`
			ScanType                     string  `json:"scanType"`
			SchemaRevision               int     `json:"schemaRevision"`
		} `json:"mediaInfo"`
		ID int `json:"id"`
	} `json:"movieFile"`
	QualityProfileID int `json:"qualityProfileId"`
	ID               int `json:"id"`
}

// Quality movie quality
type Quality struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Source     string `json:"source"`
	Resolution string `json:"resolution"`
	Modifier   string `json:"modifier"`
}

// Movies multiple Radarr movies
type Movies []*Movie

// ExcludedMovie describe an excluded movie from being downloaded
type ExcludedMovie struct {
	ID         int    `json:"id"`
	MovieTitle string `json:"movieTitle"`
	MovieYear  int    `json:"movieYear"`
	TmdbID     int    `json:"tmdbId"`
}

// ExcludedMovies describes a set of excluded movies
type ExcludedMovies []*ExcludedMovie

// DeleteMovieOptions optional option while deleting movie
type DeleteMovieOptions struct {
	// If true the movie folder and all files will be deleted when the movie is deleted
	DeleteFiles bool

	// If true the movie TMDB ID will be added to the import exclusions list when the movie is deleted
	AddExclusion bool
}

// MovieService contains Radarr movies operations
type MovieService struct {
	s      *Service
	Lookup *LookupService
}

func newMovieService(s *Service) *MovieService {
	return &MovieService{
		s:      s,
		Lookup: newLookupService(s),
	}
}

// UpcomingOptions describe period to search upcoming movies with
type UpcomingOptions struct {
	Start *time.Time
	End   *time.Time
}

// Get Returns all Movies in your collection
// https://github.com/Radarr/Radarr/wiki/API:Movie#getid
func (m *MovieService) Get(movieID int) (*Movie, error) {
	movieURL := fmt.Sprintf("%s/api%s/%d", m.s.url, movieURI, movieID)
	resp, err := m.s.client.Get(movieURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	err = parseRadarrResponse(resp)
	if err != nil {
		return nil, err
	}

	var movie Movie
	err = json.NewDecoder(resp.Body).Decode(&movie)
	if err != nil {
		return nil, err
	}

	return &movie, nil
}

// List Returns the movie with the matching ID or eerror if no matching movie is found
// https://github.com/Radarr/Radarr/wiki/API:Movie#get
func (m *MovieService) List() (Movies, error) {
	moviesURL := fmt.Sprintf("%s/api%s", m.s.url, movieURI)
	resp, err := m.s.client.Get(moviesURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	err = parseRadarrResponse(resp)
	if err != nil {
		return nil, err
	}

	var movies Movies
	err = json.NewDecoder(resp.Body).Decode(&movies)
	if err != nil {
		return nil, err
	}

	return movies, nil
}

// Upcoming Gets upcoming movies from your Radarr library, if start/end are not supplied movies airing today and tomorrow will be returned
// Its match the physicalRelease attribute
// https://github.com/Radarr/Radarr/wiki/API:Calendar#get
func (m *MovieService) Upcoming(opts ...*UpcomingOptions) (Movies, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/api%s", m.s.url, upcomingURI), nil)
	if err != nil {
		return nil, err
	}

	// If option is provided, incule them in the request
	if len(opts) > 0 {

		// If both dates are filled, verify order
		if opts[0].Start != nil && opts[0].End != nil {
			if opts[0].End.Before(*opts[0].Start) || opts[0].Start.After(*opts[0].End) {
				return nil, errors.New("Incorrect dates. Please ensure date are set properly")
			}
		}

		params := req.URL.Query()

		// If start date is defined
		if opts[0].Start != nil {
			params.Add("start", opts[0].Start.Format(time.RFC3339))
			req.URL.RawQuery = params.Encode()
		}

		// If end date is defined
		if opts[0].End != nil {
			params.Add("end", opts[0].End.Format(time.RFC3339))
			req.URL.RawQuery = params.Encode()
		}
	}

	resp, err := m.s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	err = parseRadarrResponse(resp)
	if err != nil {
		return nil, err
	}

	var movies Movies
	err = json.NewDecoder(resp.Body).Decode(&movies)
	if err != nil {
		return nil, err
	}

	return movies, nil
}

// Delete given movie
// https://github.com/Radarr/Radarr/wiki/API:Movie#deleteid
func (m *MovieService) Delete(movie *Movie, opts ...*DeleteMovieOptions) error {
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/api%s/%d", m.s.url, movieURI, movie.ID), nil)
	if err != nil {
		return err
	}

	// If option given, parse and send to request
	if len(opts) > 0 {
		d := opts[0]
		params := req.URL.Query()
		params.Add("deleteFiles", strconv.FormatBool(d.DeleteFiles))
		params.Add("addExclusion", strconv.FormatBool(d.AddExclusion))
		req.URL.RawQuery = params.Encode()
	}

	resp, err := m.s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	err = parseRadarrResponse(resp)
	return err
}

// Excluded Gets movies marked as List Exclusions
// https://github.com/Radarr/Radarr/wiki/API:List-Exclusions
func (m *MovieService) Excluded() (ExcludedMovies, error) {
	resp, err := m.s.client.Get(fmt.Sprintf("%s/api%s", m.s.url, exclusionsURI))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	err = parseRadarrResponse(resp)
	if err != nil {
		return nil, err
	}

	var movies ExcludedMovies
	err = json.NewDecoder(resp.Body).Decode(&movies)
	if err != nil {
		return nil, err
	}

	return movies, nil
}

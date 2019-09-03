package pixiv

import (
	"io"
	"net/http"
	"time"
)

// MetaSinglePage info
type MetaSinglePage struct {
	OriginalImageURL string `json:"original_image_url"`
}

// MetaPage info
type MetaPage struct {
	ImageURLs map[ImageSize]string
}

// ImageSize option
type ImageSize string

// The different ImageSizes
const (
	SizeOriginal     ImageSize = "original"
	SizeMedium                 = "medium"
	SizeLarge                  = "large"
	SizeSquareMedium           = "square_medium"
)

// SearchTarget types
type SearchTarget string

// The different search targets
const (
	SearchPartialTags     SearchTarget = "partial_match_for_tags"
	SearchExactTags                    = "exact_match_for_tags"
	SearchTitleAndCaption              = "title_and_caption"
)

// Illust Details
type Illust struct {
	MetaPage

	ID    int    `json:"id"`
	Title string `json:"title"`
	// Type  string `json:"type"`
	Caption        string
	Restrict       int
	User           Account        `json:"user"`
	Tags           Tags           `json:"tags"`
	Tools          []string       `json:"tools"`
	CreateDate     time.Time      `json:"create_date"`
	PageCount      int            `json:"page_count"`
	Width          int            `json:"width"`
	Height         int            `json:"height"`
	SanityLevel    int            `json:"sanity_level"`
	XRestrict      int            `json:"x_restrict"`
	Series         interface{}    `json:"series"` // FIXME: Unkown format
	MetaSinglePage MetaSinglePage `json:"meta_single_page"`
	MetaPages      []MetaPage     `json:"meta_pages"`
	TotalView      int            `json:"total_view"`
	TotalBookmarks int            `json:"total_bookmarks"`
	TotalComments  int            `json:"total_comments"`
	IsBookmarked   bool           `json:"is_bookmarked"`
	Visible        bool           `json:"visible"`
	IsMuted        bool           `json:"is_muted"`
}

// Illusts list
type Illusts struct {
	Illusts []Illust `json:"illusts,omitempty"`
	/// ranking_illusts
	ContextExists bool `json:"context_exists,omitempty"`
	// privacy_policy
	// next_url
}

// IllustRecommended requests your recomennded illustrations
func (s *Session) IllustRecommended(page int) (*Illusts, error) {
	if err := s.refreshAuth(); err != nil {
		return nil, err
	}

	res := new(Illusts)
	apiErr := new(APIError)

	offset := (page - 1) * 30
	if offset < 0 {
		offset = 0
	}

	_, err := s.r.New().Get("/v1/illust/recommended").QueryStruct(struct {
		Offset int `url:"offset"`
	}{
		Offset: offset,
	}).Receive(&res, &apiErr)
	if err != nil {
		return nil, err
	}

	if apiErr.HasError() {
		return nil, apiErr
	}

	return res, nil
}

// IllustDetail get a detailed illust with id
func (s *Session) IllustDetail(id uint64) (*Illust, error) {
	if err := s.refreshAuth(); err != nil {
		return nil, err
	}

	res := struct {
		Illust *Illust `json:"illust"`
	}{}

	apiErr := new(APIError)

	_, err := s.r.New().Get("v1/illust/detail").QueryStruct(struct {
		IllustID uint64 `url:"illust_id"`
	}{
		IllustID: id,
	}).Receive(&res, &apiErr)
	if err != nil {
		return nil, err
	}

	if apiErr.HasError() {
		return nil, apiErr
	}

	return res.Illust, nil
}

// IllustSearch searches for illustrations
func (s *Session) IllustSearch(term string, target SearchTarget, page int) (*Illusts, error) {
	if err := s.refreshAuth(); err != nil {
		return nil, err
	}

	res := new(Illusts)
	apiErr := new(APIError)

	offset := (page - 1) * 30
	if offset < 0 {
		offset = 0
	}

	_, err := s.r.New().Get("v1/search/illust").QueryStruct(struct {
		Term   string       `url:"word"`
		Target SearchTarget `url:"search_target"`
		Offset int          `url:"offset"`
	}{
		Term:   term,
		Target: target,
		Offset: offset,
	}).Receive(&res, &apiErr)
	if err != nil {
		return nil, err
	}

	if apiErr.HasError() {
		return nil, apiErr
	}

	return res, nil
}

// DownloadLink fetches the download link of the illustration
func (s *Session) DownloadLink(i *Illust, size ImageSize, page int) string {
	uri := i.ImageURLs[size]
	if len(i.MetaPages) > 0 {
		uri = i.MetaPages[page].ImageURLs[size]
	} else if size == SizeOriginal {
		uri = i.MetaSinglePage.OriginalImageURL
	}

	return uri
}

// Download an image from the illustration
//
// If no meta pages are set, page parameter is ignored
func (s *Session) Download(i *Illust, size ImageSize, page int) (io.ReadCloser, error) {
	uri := s.DownloadLink(i, size, page)

	r := s.r.New()
	req, err := r.Set("Referer", "https://app-api.pixiv.net/").Get(uri).Request()
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	return res.Body, nil
}

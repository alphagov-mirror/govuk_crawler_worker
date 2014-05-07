package govuk_crawler_worker_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	. "github.com/alphagov/govuk_crawler_worker"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func testServer(status int, body string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		fmt.Fprintln(w, body)
	}))
}

var _ = Describe("Crawl", func() {
	var (
		crawler    *Crawler
		crawlerErr error
	)

	BeforeEach(func() {
		crawler, crawlerErr = NewCrawler("http://127.0.0.1/")

		Expect(crawlerErr).To(BeNil())
		Expect(crawler).ToNot(BeNil())
	})

	Describe("RetryStatusCodes", func() {
		It("should return a fixed int array with values 429, 500..599", func() {
			statusCodes := RetryStatusCodes()

			Expect(len(statusCodes)).To(Equal(101))
			Expect(statusCodes[0]).To(Equal(429))
			Expect(statusCodes[1]).To(Equal(500))
			Expect(statusCodes[100]).To(Equal(599))
		})
	})

	Describe("NewCrawler()", func() {
		It("doesn't allow providing empty URLs", func() {
			badCrawler, err := NewCrawler("")

			Expect(err).ToNot(BeNil())
			Expect(badCrawler).To(BeNil())
		})

		It("provides a new crawler that accepts the provided host", func() {
			GOVUKCrawler, err := NewCrawler("https://www.gov.uk/")

			Expect(err).To(BeNil())
			Expect(GOVUKCrawler.RootURL.Host).To(Equal("www.gov.uk"))
		})
	})

	Describe("Crawler.Crawl()", func() {
		It("returns a body with no errors for 200 OK responses", func() {
			ts := testServer(http.StatusOK, "Hello world")
			defer ts.Close()

			body, err := crawler.Crawl(ts.URL)

			Expect(err).To(BeNil())
			Expect(strings.TrimSpace(string(body))).To(Equal("Hello world"))
		})

		It("doesn't allow crawling a URL that doesn't match the root URL", func() {
			body, err := crawler.Crawl("http://google.com/foo")

			Expect(err).To(Equal(CannotCrawlURL))
			Expect(body).To(Equal([]byte{}))
		})

		Describe("returning a retry error", func() {
			It("returns a retry error if we get a response code of Too Many Requests", func() {
				ts := testServer(429, "Too Many Requests")
				defer ts.Close()

				body, err := crawler.Crawl(ts.URL)

				Expect(err).To(Equal(RetryRequestError))
				Expect(body).To(Equal([]byte{}))
			})

			It("returns a retry error if we get a response code of Internal Server Error", func() {
				ts := testServer(http.StatusInternalServerError, "Internal Server Error")
				defer ts.Close()

				body, err := crawler.Crawl(ts.URL)

				Expect(err).To(Equal(RetryRequestError))
				Expect(body).To(Equal([]byte{}))
			})

			It("returns a retry error if we get a response code of Gateway Timeout", func() {
				ts := testServer(http.StatusGatewayTimeout, "Gateway Timeout")
				defer ts.Close()

				body, err := crawler.Crawl(ts.URL)

				Expect(err).To(Equal(RetryRequestError))
				Expect(body).To(Equal([]byte{}))
			})
		})
	})
})
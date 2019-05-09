package orderS

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestValidateExternalURL(t *testing.T) {
	Convey("example.com", t, func() {
		recognizedURLs := []string{"example.com"}

		Convey("Valid", func() {
			err := validateExternalURL(recognizedURLs, "https://example.com/123")
			So(err, ShouldBeNil)
		})
		Convey("Invalid", func() {
			err := validateExternalURL(recognizedURLs, "https://hello.example.com/123")
			So(err, ShouldNotBeNil)
		})
		Convey("Only allow http and https", func() {
			err := validateExternalURL(recognizedURLs, "app://hello.example.com/123")
			So(err, ShouldNotBeNil)
		})
		Convey("Only allow http and https (2)", func() {
			err := validateExternalURL(recognizedURLs, "//hello.example.com/123")
			So(err, ShouldNotBeNil)
		})
	})

	Convey("*.example.com", t, func() {
		recognizedURLs := []string{"*.example.com"}

		Convey("Valid", func() {
			err := validateExternalURL(recognizedURLs, "https://hello-world.example.com/123")
			So(err, ShouldBeNil)
		})
		Convey("Wildcard (invalid)", func() {
			err := validateExternalURL(recognizedURLs, "https://*.example.com/123")
			So(err, ShouldNotBeNil)
		})
		Convey("Multiple level of subdomain (invalid)", func() {
			err := validateExternalURL(recognizedURLs, "https://abc.foo.example.com/123")
			So(err, ShouldNotBeNil)
		})
		Convey("Start with a dot (invalid)", func() {
			err := validateExternalURL(recognizedURLs, "https://.example.com/123")
			So(err, ShouldNotBeNil)
		})
		Convey("No subdomain (invalid)", func() {
			err := validateExternalURL(recognizedURLs, "https://example.com/123")
			So(err, ShouldNotBeNil)
		})
	})
}

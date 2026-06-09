package providers_test

import (
	"errors"

	"github.com/m4rc3l05/media-follower/internal/providers"
	"github.com/m4rc3l05/media-follower/internal/providers/inputs"
	"github.com/m4rc3l05/media-follower/internal/providers/outputs"
	"github.com/m4rc3l05/media-follower/internal/test"
	"github.com/m4rc3l05/media-follower/internal/test/testdata"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ItunesMusicReleasesProvider", func() {
	Describe("FetchReleases()", func() {
		It("return error if output provider returns error", func() {
			err := errors.New("foo")
			in := inputs.MockInputProvider[inputs.ItunesArtist]{
				InterfaceMock: test.InterfaceMock{
					Calls:       map[string][]test.CallInfo{},
					MockReturns: map[string][][]any{},
				},
			}
			out := outputs.MockOutputProvider[inputs.ItunesArtist, outputs.ItunesAlbum]{
				InterfaceMock: test.InterfaceMock{
					Calls: map[string][]test.CallInfo{},
					MockReturns: map[string][][]any{
						"FetchOutputs": {{[]outputs.ItunesAlbum{}, err}},
					},
				},
			}

			provider := providers.NewItunesMusicReleasesProvider(in, out)

			Expect(provider.FetchReleases(inputs.ItunesArtist{})).Error().To(MatchError(err))
		})

		It("return only return valid releases", func() {
			ok := testdata.OkItunesAlbumStruct()
			in := inputs.MockInputProvider[inputs.ItunesArtist]{
				InterfaceMock: test.InterfaceMock{
					Calls:       map[string][]test.CallInfo{},
					MockReturns: map[string][][]any{},
				},
			}
			out := outputs.MockOutputProvider[inputs.ItunesArtist, outputs.ItunesAlbum]{
				InterfaceMock: test.InterfaceMock{
					Calls: map[string][]test.CallInfo{},
					MockReturns: map[string][][]any{
						"FetchOutputs": {{[]outputs.ItunesAlbum{
							ok,
							testdata.BadItunesAlbumNoReleaseStruct(),
							testdata.BadItunesAlbumDJMixStruct(),
							testdata.BadItunesAlbumDJMix2Struct(),
							testdata.BadItunesAlbumCompilationStruct(),
						}, nil}},
					},
				},
			}

			provider := providers.NewItunesMusicReleasesProvider(in, out)
			releases, err := provider.FetchReleases(inputs.ItunesArtist{})

			Expect(err).To(BeNil())
			Expect(releases).To(HaveLen(1))
			Expect(releases[0]).To(Equal(ok))
		})
	})
})

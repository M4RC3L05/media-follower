package factories

import (
	"fmt"

	"github.com/m4rc3l05/media-follower/internal/common/providers"
	"github.com/m4rc3l05/media-follower/internal/common/utils"
)

func ProviderFactory(name providers.ProviderName) providers.IReleaseProvider {
	if name == providers.ITUNES_MUSIC_RELEASES_PROVIDER {
		return providers.ItunesMusicReleaseProvider{
			Validator: utils.NewValidator(),
			Conform:   utils.NewModifier(),
		}
	}

	panic(fmt.Sprintf("Provider \"%s\" not supported", name))
}

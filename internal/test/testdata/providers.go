package testdata

import (
	"github.com/m4rc3l05/media-follower/internal/providers"
)

func OkItunesAlbumStruct() providers.ItunesAlbum {
	m := OkItunesAlbumHttpResponse()
	return providers.ItunesAlbum{
		AmgArtistID:            new(m["amgArtistId"].(int64)),
		ArtistID:               m["artistId"].(int64),
		ArtistName:             m["artistName"].(string),
		ArtistViewURL:          m["artistViewUrl"].(string),
		ArtworkURL100:          m["artworkUrl100"].(string),
		ArtworkURL60:           m["artworkUrl60"].(string),
		CollectionCensoredName: m["collectionCensoredName"].(string),
		CollectionExplicitness: m["collectionExplicitness"].(string),
		CollectionID:           m["collectionId"].(int64),
		CollectionName:         m["collectionName"].(string),
		CollectionPrice:        new(m["collectionPrice"].(float64)),
		CollectionType:         m["collectionType"].(string),
		CollectionViewURL:      m["collectionViewUrl"].(string),
		ContentAdvisoryRating:  new(m["contentAdvisoryRating"].(string)),
		Copyright:              new(m["copyright"].(string)),
		Country:                m["country"].(string),
		Currency:               m["currency"].(string),
		PrimaryGenreName:       m["primaryGenreName"].(string),
		ReleaseDate:            new(m["releaseDate"].(string)),
		TrackCount:             m["trackCount"].(int64),
		WrapperType:            m["wrapperType"].(string),
	}
}

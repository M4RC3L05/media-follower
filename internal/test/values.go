package test

import "github.com/m4rc3l05/media-follower/internal/providers"

func OkItunesAlbumMap() map[string]any {
	return map[string]any{
		"amgArtistId":            int64(1234567890),
		"artistId":               int64(987654321),
		"artistName":             "Example Artist",
		"artistViewUrl":          "https://music.example.com/artist/987654321",
		"artworkUrl100":          "https://images.example.com/artwork100.jpg",
		"artworkUrl60":           "https://images.example.com/artwork60.jpg",
		"collectionCensoredName": "Example Album (Clean)",
		"collectionExplicitness": "notExplicit",
		"collectionId":           int64(555666777),
		"collectionName":         "Example Album",
		"collectionPrice":        float64(9.99),
		"collectionType":         "Album",
		"collectionViewUrl":      "https://music.example.com/album/555666777",
		"contentAdvisoryRating":  "PG",
		"copyright":              "© 2020 Example Records",
		"country":                "USA",
		"currency":               "USD",
		"primaryGenreName":       "Pop",
		"releaseDate":            "2026-06-05T14:23:45Z",
		"trackCount":             int64(12),
		"wrapperType":            "collection",
	}
}

func OkItunesAlbumStruct() providers.ItunesAlbum {
	m := OkItunesAlbumMap()
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

func BadItunesAlbumCompilationMap() map[string]any {
	return map[string]any{
		"amgArtistId":            int64(1234567890),
		"artistId":               int64(987654321),
		"artistName":             "Various Artists",
		"artistViewUrl":          "https://music.example.com/artist/987654321",
		"artworkUrl100":          "https://images.example.com/artwork100.jpg",
		"artworkUrl60":           "https://images.example.com/artwork60.jpg",
		"collectionCensoredName": "Example Album (Clean)",
		"collectionExplicitness": "notExplicit",
		"collectionId":           int64(247358353),
		"collectionName":         "Example Album",
		"collectionPrice":        float64(9.99),
		"collectionType":         "Album",
		"collectionViewUrl":      "https://music.example.com/album/555666777",
		"contentAdvisoryRating":  "PG",
		"copyright":              "© 2020 Example Records",
		"country":                "USA",
		"currency":               "USD",
		"primaryGenreName":       "Pop",
		"releaseDate":            "2026-06-05T14:23:45Z",
		"trackCount":             int64(12),
		"wrapperType":            "collection",
	}
}

func BadItunesAlbumNoReleaseMap() map[string]any {
	return map[string]any{
		"amgArtistId":            int64(1234567890),
		"artistId":               int64(987654321),
		"artistName":             "Various Artists",
		"artistViewUrl":          "https://music.example.com/artist/987654321",
		"artworkUrl100":          "https://images.example.com/artwork100.jpg",
		"artworkUrl60":           "https://images.example.com/artwork60.jpg",
		"collectionCensoredName": "Example Album (Clean)",
		"collectionExplicitness": "notExplicit",
		"collectionId":           int64(55565867666777),
		"collectionName":         "Various Artists",
		"collectionPrice":        float64(9.99),
		"collectionType":         "Album",
		"collectionViewUrl":      "https://music.example.com/album/555666777",
		"contentAdvisoryRating":  "PG",
		"copyright":              "© 2020 Example Records",
		"country":                "USA",
		"currency":               "USD",
		"primaryGenreName":       "Pop",
		"trackCount":             int64(12),
		"wrapperType":            "collection",
	}
}

func BadItunesAlbumDJMixMap() map[string]any {
	return map[string]any{
		"amgArtistId":            int64(1234567890),
		"artistId":               int64(987654321),
		"artistName":             "Various Artists",
		"artistViewUrl":          "https://music.example.com/artist/987654321",
		"artworkUrl100":          "https://images.example.com/artwork100.jpg",
		"artworkUrl60":           "https://images.example.com/artwork60.jpg",
		"collectionCensoredName": "Example Album (Clean)",
		"collectionExplicitness": "notExplicit",
		"collectionId":           int64(79658657),
		"collectionName":         "Various Artists",
		"collectionPrice":        float64(9.99),
		"collectionType":         "Album",
		"collectionViewUrl":      "https://music.example.com/album/555666777",
		"contentAdvisoryRating":  "PG",
		"copyright":              "© 2020 Example Records",
		"country":                "USA",
		"currency":               "USD",
		"primaryGenreName":       "Pop",
		"releaseDate":            "2026-06-05T14:23:45Z",
		"trackCount":             int64(12),
		"wrapperType":            "collection",
	}
}

func BadItunesAlbumDJMix2Map() map[string]any {
	return map[string]any{
		"amgArtistId":            int64(1234567890),
		"artistId":               int64(987654321),
		"artistName":             "Various Artists",
		"artistViewUrl":          "https://music.example.com/artist/987654321",
		"artworkUrl100":          "https://images.example.com/artwork100.jpg",
		"artworkUrl60":           "https://images.example.com/artwork60.jpg",
		"collectionCensoredName": "Various Artists (Clean)",
		"collectionExplicitness": "notExplicit",
		"collectionId":           int64(7436227447),
		"collectionName":         "Example Album",
		"collectionPrice":        float64(9.99),
		"collectionType":         "Album",
		"collectionViewUrl":      "https://music.example.com/album/555666777",
		"contentAdvisoryRating":  "PG",
		"copyright":              "© 2020 Example Records",
		"country":                "USA",
		"currency":               "USD",
		"primaryGenreName":       "Pop",
		"releaseDate":            "2026-06-05T14:23:45Z",
		"trackCount":             int64(12),
		"wrapperType":            "collection",
	}
}

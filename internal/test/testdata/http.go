package testdata

func OkItunesAlbumHttpResponse() map[string]any {
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

func BadItunesAlbumCompilationHttpResponse() map[string]any {
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

func BadItunesAlbumNoReleaseHttpResponse() map[string]any {
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

func BadItunesAlbumDJMixHttpResponse() map[string]any {
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

func BadItunesAlbumDJMix2HttpResponse() map[string]any {
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

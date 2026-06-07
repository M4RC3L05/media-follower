package passwordhashing_test

import (
	"strconv"
	"strings"

	passwordhashing "github.com/m4rc3l05/media-follower/internal/common/password_hashing"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Argon2di", func() {
	Describe("Hash()", func() {
		It("should hash a plaintex", func() {
			pm := passwordhashing.NewArgon2di()
			pm.Memory = 1
			pm.Threads = 1
			pm.Time = 1
			stringHash := pm.Hash("foo")
			stringHash2 := pm.Hash("foo")
			parts := strings.Split(stringHash, " ")

			Expect(stringHash).NotTo(Equal(stringHash2))
			Expect(parts).To(HaveLen(6))
			Expect(parts[0]).To(Equal(strconv.FormatUint(uint64(pm.Time), 10)))
			Expect(parts[1]).To(Equal(strconv.FormatUint(uint64(pm.Memory), 10)))
			Expect(parts[2]).To(Equal(strconv.FormatUint(uint64(pm.Threads), 10)))
			Expect(parts[3]).To(Equal(strconv.FormatUint(uint64(pm.Keylen), 10)))
		})
	})

	Describe("Compare()", func() {
		It("should return false if plaintex and hash do not match", func() {
			pm := passwordhashing.NewArgon2di()
			pm.Memory = 1
			pm.Threads = 1
			pm.Time = 1
			stringHash := pm.Hash("foo")

			Expect(pm.Compare(stringHash, "bar")).To(BeFalse())
		})

		It("should return true if plaintex and hash match", func() {
			pm := passwordhashing.NewArgon2di()
			pm.Memory = 1
			pm.Threads = 1
			pm.Time = 1
			stringHash := pm.Hash("foo")

			Expect(pm.Compare(stringHash, "foo")).To(BeTrue())
		})

		It(
			"should return true if plaintex and hash match with hash encoded with different params",
			func() {
				pm := passwordhashing.NewArgon2di()
				pm.Memory = 2
				pm.Threads = 1
				pm.Time = 1

				Expect(
					pm.Compare(
						"1 1 1 32 +SAS6AuVDgQokpurmfU5eg zQyjVEgstmlfp2SzgjISm54oHTeVlZZKjF7W8p0TUSA",
						"foo",
					),
				).To(BeTrue())
			},
		)
	})
})

package tests

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/Velocidex/ordereddict"
	"github.com/alecthomas/assert"
	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/suite"
)

type FATTestSuite struct {
	suite.Suite
	binary, extension string
	tmpdir            string
}

func (self *FATTestSuite) SetupTest() {
	if runtime.GOOS == "windows" {
		self.extension = ".exe"
	}

	// Search for a valid binary to run.
	binaries, err := filepath.Glob(
		"../fat" + self.extension)
	assert.NoError(self.T(), err)

	self.binary, _ = filepath.Abs(binaries[0])
	fmt.Printf("Found binary %v\n", self.binary)

	self.tmpdir, err = ioutil.TempDir("", "tmp")
	assert.NoError(self.T(), err)
}

func (self *FATTestSuite) TearDownTest() {
	os.RemoveAll(self.tmpdir)
}

func (self *FATTestSuite) TestFAT12Support() {
	cmd := exec.Command(self.binary, "ls", "../testcases/fat12.dd", "a/b")
	out, err := cmd.CombinedOutput()
	assert.NoError(self.T(), err, string(out))

	golden := ordereddict.NewDict().Set("listdir", strings.Split(string(out), "\n"))

	cmd = exec.Command(self.binary, "cat", "../testcases/fat12.dd", "a/b/alice.txt")
	out, err = cmd.CombinedOutput()
	assert.NoError(self.T(), err, string(out))

	md5_sum := md5.New()
	md5_sum.Write(out)

	golden.Set("alice_md5", hex.EncodeToString(md5_sum.Sum(nil)))

	cmd = exec.Command(self.binary, "stat", "../testcases/fat12.dd", "a/b/alice.txt")
	out, err = cmd.CombinedOutput()
	assert.NoError(self.T(), err, string(out))

	golden.Set("stat", strings.Split(string(out), "\n"))

	g := goldie.New(self.T(), goldie.WithFixtureDir("fixtures"))
	g.AssertJson(self.T(), "fat12", golden)
}

func (self *FATTestSuite) TestFAT16Support() {
	cmd := exec.Command(self.binary, "ls", "../testcases/fat16.dd", "a/b")
	out, err := cmd.CombinedOutput()
	assert.NoError(self.T(), err, string(out))

	golden := ordereddict.NewDict().Set("listdir", strings.Split(string(out), "\n"))

	cmd = exec.Command(self.binary, "cat", "../testcases/fat16.dd", "a/b/alice.txt")
	out, err = cmd.CombinedOutput()
	assert.NoError(self.T(), err, string(out))

	md5_sum := md5.New()
	md5_sum.Write(out)

	golden.Set("alice_md5", hex.EncodeToString(md5_sum.Sum(nil)))

	cmd = exec.Command(self.binary, "stat", "../testcases/fat16.dd", "a/b/alice.txt")
	out, err = cmd.CombinedOutput()
	assert.NoError(self.T(), err, string(out))

	golden.Set("stat", strings.Split(string(out), "\n"))

	g := goldie.New(self.T(), goldie.WithFixtureDir("fixtures"))
	g.AssertJson(self.T(), "fat16", golden)
}

func (self *FATTestSuite) TestFAT32Support() {
	cmd := exec.Command(self.binary, "ls", "../testcases/fat32.dd", "a/b")
	out, err := cmd.CombinedOutput()
	assert.NoError(self.T(), err, string(out))

	golden := ordereddict.NewDict().Set("listdir", strings.Split(string(out), "\n"))

	cmd = exec.Command(self.binary, "cat", "../testcases/fat32.dd", "a/b/alice.txt")
	out, err = cmd.CombinedOutput()
	assert.NoError(self.T(), err, string(out))

	md5_sum := md5.New()
	md5_sum.Write(out)

	golden.Set("alice_md5", hex.EncodeToString(md5_sum.Sum(nil)))

	cmd = exec.Command(self.binary, "stat", "../testcases/fat32.dd", "a/b/alice.txt")
	out, err = cmd.CombinedOutput()
	assert.NoError(self.T(), err, string(out))

	golden.Set("stat", strings.Split(string(out), "\n"))

	// List a very large directory
	cmd = exec.Command(self.binary, "ls", "../testcases/fat32.dd", "a/b/c")
	out, err = cmd.CombinedOutput()
	assert.NoError(self.T(), err, string(out))

	filtered := []string{}
	for _, l := range strings.Split(string(out), "\n") {
		if strings.Contains(l, "a_long_filename_with_extra_characters") {
			filtered = append(filtered, l)
		}
	}

	golden.Set("listdir_large", filtered)

	g := goldie.New(self.T(), goldie.WithFixtureDir("fixtures"))
	g.AssertJson(self.T(), "fat32", golden)
}

func TestFAT(t *testing.T) {
	suite.Run(t, &FATTestSuite{})
}

func formatForTest(v interface{}) []string {
	serialized, _ := json.MarshalIndent(v, " ", " ")
	return strings.Split(string(serialized), "\n")
}

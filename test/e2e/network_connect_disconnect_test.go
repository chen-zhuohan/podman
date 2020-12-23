package integration

import (
	"os"

	. "github.com/containers/podman/v2/test/utils"
	"github.com/containers/storage/pkg/stringid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Podman network connect and disconnect", func() {
	var (
		tempdir    string
		err        error
		podmanTest *PodmanTestIntegration
	)

	BeforeEach(func() {
		tempdir, err = CreateTempDirInTempDir()
		if err != nil {
			os.Exit(1)
		}
		podmanTest = PodmanTestCreate(tempdir)
		podmanTest.Setup()
	})

	AfterEach(func() {
		podmanTest.Cleanup()
		f := CurrentGinkgoTestDescription()
		processTestResult(f)

	})

	It("bad network name in disconnect should result in error", func() {
		SkipIfRootless("network connect and disconnect are only rootfull")
		dis := podmanTest.Podman([]string{"network", "disconnect", "foobar", "test"})
		dis.WaitWithDefaultTimeout()
		Expect(dis.ExitCode()).ToNot(BeZero())

	})

	It("bad container name in network disconnect should result in error", func() {
		SkipIfRootless("network connect and disconnect are only rootfull")
		netName := "aliasTest" + stringid.GenerateNonCryptoID()
		session := podmanTest.Podman([]string{"network", "create", netName})
		session.WaitWithDefaultTimeout()
		Expect(session.ExitCode()).To(BeZero())
		defer podmanTest.removeCNINetwork(netName)

		dis := podmanTest.Podman([]string{"network", "disconnect", netName, "foobar"})
		dis.WaitWithDefaultTimeout()
		Expect(dis.ExitCode()).ToNot(BeZero())

	})

	It("podman network disconnect", func() {
		SkipIfRootless("network connect and disconnect are only rootfull")
		netName := "aliasTest" + stringid.GenerateNonCryptoID()
		session := podmanTest.Podman([]string{"network", "create", netName})
		session.WaitWithDefaultTimeout()
		Expect(session.ExitCode()).To(BeZero())
		defer podmanTest.removeCNINetwork(netName)

		ctr := podmanTest.Podman([]string{"run", "-dt", "--name", "test", "--network", netName, ALPINE, "top"})
		ctr.WaitWithDefaultTimeout()
		Expect(ctr.ExitCode()).To(BeZero())

		exec := podmanTest.Podman([]string{"exec", "-it", "test", "ip", "addr", "show", "eth0"})
		exec.WaitWithDefaultTimeout()
		Expect(exec.ExitCode()).To(BeZero())

		dis := podmanTest.Podman([]string{"network", "disconnect", netName, "test"})
		dis.WaitWithDefaultTimeout()
		Expect(dis.ExitCode()).To(BeZero())

		exec = podmanTest.Podman([]string{"exec", "-it", "test", "ip", "addr", "show", "eth0"})
		exec.WaitWithDefaultTimeout()
		Expect(exec.ExitCode()).ToNot(BeZero())
	})

	It("bad network name in connect should result in error", func() {
		SkipIfRootless("network connect and disconnect are only rootfull")
		dis := podmanTest.Podman([]string{"network", "connect", "foobar", "test"})
		dis.WaitWithDefaultTimeout()
		Expect(dis.ExitCode()).ToNot(BeZero())

	})

	It("bad container name in network connect should result in error", func() {
		SkipIfRootless("network connect and disconnect are only rootfull")
		netName := "aliasTest" + stringid.GenerateNonCryptoID()
		session := podmanTest.Podman([]string{"network", "create", netName})
		session.WaitWithDefaultTimeout()
		Expect(session.ExitCode()).To(BeZero())
		defer podmanTest.removeCNINetwork(netName)

		dis := podmanTest.Podman([]string{"network", "connect", netName, "foobar"})
		dis.WaitWithDefaultTimeout()
		Expect(dis.ExitCode()).ToNot(BeZero())

	})

	It("podman connect on a container that already is connected to the network should error", func() {
		SkipIfRootless("network connect and disconnect are only rootfull")
		netName := "aliasTest" + stringid.GenerateNonCryptoID()
		session := podmanTest.Podman([]string{"network", "create", netName})
		session.WaitWithDefaultTimeout()
		Expect(session.ExitCode()).To(BeZero())
		defer podmanTest.removeCNINetwork(netName)

		ctr := podmanTest.Podman([]string{"create", "--name", "test", "--network", netName, ALPINE, "top"})
		ctr.WaitWithDefaultTimeout()
		Expect(ctr.ExitCode()).To(BeZero())

		con := podmanTest.Podman([]string{"network", "connect", netName, "test"})
		con.WaitWithDefaultTimeout()
		Expect(con.ExitCode()).ToNot(BeZero())
	})

	It("podman network connect", func() {
		SkipIfRemote("This requires a pending PR to be merged before it will work")
		SkipIfRootless("network connect and disconnect are only rootfull")
		netName := "aliasTest" + stringid.GenerateNonCryptoID()
		session := podmanTest.Podman([]string{"network", "create", netName})
		session.WaitWithDefaultTimeout()
		Expect(session.ExitCode()).To(BeZero())
		defer podmanTest.removeCNINetwork(netName)

		ctr := podmanTest.Podman([]string{"run", "-dt", "--name", "test", "--network", netName, ALPINE, "top"})
		ctr.WaitWithDefaultTimeout()
		Expect(ctr.ExitCode()).To(BeZero())

		exec := podmanTest.Podman([]string{"exec", "-it", "test", "ip", "addr", "show", "eth0"})
		exec.WaitWithDefaultTimeout()
		Expect(exec.ExitCode()).To(BeZero())

		// Create a second network
		newNetName := "aliasTest" + stringid.GenerateNonCryptoID()
		session = podmanTest.Podman([]string{"network", "create", newNetName})
		session.WaitWithDefaultTimeout()
		Expect(session.ExitCode()).To(BeZero())
		defer podmanTest.removeCNINetwork(newNetName)

		connect := podmanTest.Podman([]string{"network", "connect", newNetName, "test"})
		connect.WaitWithDefaultTimeout()
		Expect(connect.ExitCode()).To(BeZero())

		exec = podmanTest.Podman([]string{"exec", "-it", "test", "ip", "addr", "show", "eth1"})
		exec.WaitWithDefaultTimeout()
		Expect(exec.ExitCode()).To(BeZero())
	})

	It("podman network connect when not running", func() {
		SkipIfRootless("network connect and disconnect are only rootfull")
		netName := "aliasTest" + stringid.GenerateNonCryptoID()
		session := podmanTest.Podman([]string{"network", "create", netName})
		session.WaitWithDefaultTimeout()
		Expect(session.ExitCode()).To(BeZero())
		defer podmanTest.removeCNINetwork(netName)

		ctr := podmanTest.Podman([]string{"create", "--name", "test", ALPINE, "top"})
		ctr.WaitWithDefaultTimeout()
		Expect(ctr.ExitCode()).To(BeZero())

		dis := podmanTest.Podman([]string{"network", "connect", netName, "test"})
		dis.WaitWithDefaultTimeout()
		Expect(dis.ExitCode()).To(BeZero())

		start := podmanTest.Podman([]string{"start", "test"})
		start.WaitWithDefaultTimeout()
		Expect(start.ExitCode()).To(BeZero())

		exec := podmanTest.Podman([]string{"exec", "-it", "test", "ip", "addr", "show", "eth0"})
		exec.WaitWithDefaultTimeout()
		Expect(exec.ExitCode()).To(BeZero())

		exec = podmanTest.Podman([]string{"exec", "-it", "test", "ip", "addr", "show", "eth1"})
		exec.WaitWithDefaultTimeout()
		Expect(exec.ExitCode()).To(BeZero())
	})

	It("podman network disconnect when not running", func() {
		SkipIfRootless("network connect and disconnect are only rootfull")
		netName1 := "aliasTest" + stringid.GenerateNonCryptoID()
		session := podmanTest.Podman([]string{"network", "create", netName1})
		session.WaitWithDefaultTimeout()
		Expect(session.ExitCode()).To(BeZero())
		defer podmanTest.removeCNINetwork(netName1)

		netName2 := "aliasTest" + stringid.GenerateNonCryptoID()
		session2 := podmanTest.Podman([]string{"network", "create", netName2})
		session2.WaitWithDefaultTimeout()
		Expect(session2.ExitCode()).To(BeZero())
		defer podmanTest.removeCNINetwork(netName2)

		ctr := podmanTest.Podman([]string{"create", "--name", "test", "--network", netName1 + "," + netName2, ALPINE, "top"})
		ctr.WaitWithDefaultTimeout()
		Expect(ctr.ExitCode()).To(BeZero())

		dis := podmanTest.Podman([]string{"network", "disconnect", netName1, "test"})
		dis.WaitWithDefaultTimeout()
		Expect(dis.ExitCode()).To(BeZero())

		start := podmanTest.Podman([]string{"start", "test"})
		start.WaitWithDefaultTimeout()
		Expect(start.ExitCode()).To(BeZero())

		exec := podmanTest.Podman([]string{"exec", "-it", "test", "ip", "addr", "show", "eth0"})
		exec.WaitWithDefaultTimeout()
		Expect(exec.ExitCode()).To(BeZero())

		exec = podmanTest.Podman([]string{"exec", "-it", "test", "ip", "addr", "show", "eth1"})
		exec.WaitWithDefaultTimeout()
		Expect(exec.ExitCode()).ToNot(BeZero())
	})
})

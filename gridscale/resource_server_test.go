package gridscale

import (
	"fmt"
	"testing"

	"bitbucket.org/gridscale/gsclient-go"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccDataSourceGridscaleServer_Basic(t *testing.T) {
	var object gsclient.Server
	name := fmt.Sprintf("object-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourceGridscaleServerConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceGridscaleServerExists("gridscale_server.foo", &object),
					resource.TestCheckResourceAttr(
						"gridscale_server.foo", "name", name),
					resource.TestCheckResourceAttr(
						"gridscale_server.foo", "cores", "1"),
					resource.TestCheckResourceAttr(
						"gridscale_server.foo", "memory", "1"),
				),
			},
		},
	})
}

func testAccCheckDataSourceGridscaleServerExists(n string, object *gsclient.Server) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No object UUID is set")
		}

		client := testAccProvider.Meta().(*gsclient.Client)

		id := rs.Primary.ID

		foundObject, err := client.GetServer(id)

		if err != nil {
			return err
		}

		if foundObject.Properties.ObjectUuid != id {
			return fmt.Errorf("Object not found")
		}

		*object = *foundObject

		return nil
	}
}

func testAccCheckDataSourceGridscaleServerConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "gridscale_server" "foo" {
  name   = "%s"
  cores = 1
  memory = 1
}
`, name)
}

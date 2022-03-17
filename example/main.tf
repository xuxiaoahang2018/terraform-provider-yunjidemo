terraform {
  required_providers {
    yunjidemo = {
      source  = "yunji/yunjidemo"
    }
  }
}

resource "yunjidemo_demo" "test" {
  instance_name  = "yunji"
  disk_size = "100"

}

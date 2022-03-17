terraform {
  required_providers {
    yunjidemo = {
      source  = "yunji/yunjidemo"
    }
  }
}

resource "yunjidemo_demo" "test" {
  instance_name  = "aini"
  disk_size = 100

}

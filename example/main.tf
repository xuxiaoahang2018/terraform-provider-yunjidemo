terraform {
  required_providers {
    yunjidemo = {
      source  = "yunji/yunjidemo"
    }
  }
}

resource "yunjidemo_demo" "test" {
  instance_name  = "yunji"
  disk_size = 100
  networks {
	port = 81
    protocol = "http"
  }
  networks {
	port = 82
    protocol = "https"
  }
  memory = {
    memory_size = "1024"
	memory_unit = "test"
  }
  config_json =<<EOF
{"test_json":"yunji"}
EOF
}
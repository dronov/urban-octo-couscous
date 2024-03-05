terraform {
  required_providers {
    proxmox = {
      source = "telmate/proxmox"
      version = "2.9.14"
    }
  }
}

provider "proxmox" {
  pm_api_url = "https://192.168.0.11:8006/api2/json"
  pm_api_token_id = "twtrp@pam!twtrp"
  pm_api_token_secret = "377c0076-e3ce-470f-b6aa-96f87929c801"
  pm_tls_insecure = true
}

resource "proxmox_vm_qemu" "instance" {
  count       = 5
  name        = count.index == 3 ? "infra${count.index}" : "server${count.index}"
  target_node = var.proxmox_host
  clone       = var.ubuntu_template_name
  os_type     = "linux"
  cores       = 4
  sockets     = 1
  cpu         = "host"
  memory      = count.index == 3 ? 4096 : 2048
  scsihw      = "virtio-scsi-pci"
  bootdisk    = "scsi0"

  disk {
    slot    = 0
    size    = "50G"
    type    = "scsi"
    storage = "local-lvm"
    iothread = 0
  }

  network {
    model  = "virtio"
    bridge = "vmbr0"
  }

  ipconfig0 = count.index == 3 ? "ip=192.168.0.5${count.index}/24,gw=192.168.0.1" : "ip=192.168.0.10${count.index + 1}/24,gw=192.168.0.1"

  lifecycle {
    ignore_changes = ["network"]
  }

  sshkeys = var.ssh_key
}

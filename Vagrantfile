Vagrant.configure("2") do |config|

  config.vm.box       = "ubuntu/trusty64"

  config.vm.provision "ansible" do |ansible|
    ENV["ANSIBLE_CONFIG"] = "ansible/ansible.cfg"
    ansible.playbook = "ansible/dev.yml"
    ansible.inventory_path = "ansible/dev_hosts"
    ansible.tags = ENV["TAGS"]
    ansible.skip_tags = ENV["SKIP_TAGS"]
    ansible.verbose = 'v'
    ansible.limit = "all"
  end

  config.vm.network :forwarded_port, host: 2207, guest: 22 
  config.vm.network :forwarded_port, host: 8083, guest: 8083 
  config.vm.network :forwarded_port, host: 8086, guest: 8086 

  config.vm.network "private_network", ip: "192.168.50.117"

  config.vm.synced_folder ".", "/vagrant", :disabled => true
  config.vm.synced_folder "tb-presence/", "/webapps/tb-presence", create: true, id: "vagrant-root",
    owner: "vagrant",
    group: "www-data",
    mount_options: ["dmode=775,fmode=764"]
end
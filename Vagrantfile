Vagrant.configure("2") do |config|

  config.vm.define "droplist" do |c|
  	c.vm.box = "ubuntu/xenial64"
  	c.vm.hostname = "droplist"
  	c.vm.synced_folder ".", "/data"
  	c.vm.provision "shell", inline: <<-SHELL
  		sudo su
  		
      apt install unzip -y
      
      mkdir -p /opt/droplist/
      cp -r /data/* /opt/droplist

      cd /opt/droplist/resources
      unzip droplist.zip

      echo "y" | ufw enable
      /opt/droplist/resources/droplist

  	SHELL
  end
end
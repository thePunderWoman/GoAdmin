default_run_options[:pty] = true
ssh_options[:forward_agent] = true


set :application, "GoAdmin"
set :repository,  "git@github.com:curt-labs/GoAdmin.git"

set :scm, :git
set :scm_passphrase, ""
set :user, "deployer"

#role :web, "curt-api-server1.cloudapp.net", "curt-api-server2.cloudapp.net"
role :web, "173.255.117.20", "173.255.112.170"

set :deploy_to, "/home/#{user}/#{application}"
set :deploy_via, :remote_cache

set :use_sudo, false
set :sudo_prompt, ""
set :normalize_asset_timestamps, false

#set :default_environment, {
#  'GOPATH' => "$HOME/gocode"
#}

after :deploy, "deploy:goget", "db:configure", "email:configure", "deploy:compile", "deploy:stop", "deploy:restart"

namespace :db do
  desc "set database Connection Strings"
  task :configure do
    set(:database_username) { Capistrano::CLI.ui.ask("Database Username:") }
  
    set(:database_password) { Capistrano::CLI.password_prompt("Database Password:") }

    db_config = <<-EOF
      package database

      const (
        db_proto = "tcp"
        db_addr  = "curtsql.cloudapp.net:3306"
        db_user  = "#{database_username}"
        db_pass  = "#{database_password}"
		CurtDevdb_name = "CurtDev2"
		Admindb_name   = "admin"
      )
    EOF
    run "mkdir -p #{deploy_to}/current/helpers/database"
    put db_config, "#{deploy_to}/current/helpers/database/ConnectionString.go"
  end
end
namespace :email do
  desc "set email settings"
  task :configure do
    set(:email_password) { Capistrano::CLI.password_prompt("Email Password:") }
    email_config = <<-EOF
      package email

      const (
        EmailServer   = "smtp.gmail.com"
        EmailAddress  = "no-reply@curtmfg.com"
        EmailUsername = "no-reply@curtmfg.com"
        EmailPassword = "#{email_password}"
        EmailSSL      = true
        EmailPort     = 587
      )
    EOF
    run "mkdir -p #{deploy_to}/current/helpers/email"
    put email_config, "#{deploy_to}/current/helpers/email/EmailSettings.go"
  end
end
namespace :deploy do
  task :goget do
  	run "/home/#{user}/bin/go get github.com/ziutek/mymysql/native"
    run "/home/#{user}/bin/go get github.com/ziutek/mymysql/mysql"
  	run "/home/#{user}/bin/go get github.com/gorilla/sessions"
  end
  task :compile do
  	run "GOOS=linux GOARCH=amd64 CGO_ENABLED=0 /home/#{user}/bin/go build -o #{deploy_to}/current/go-admin #{deploy_to}/current/index.go"
  end
  task :start do ; end
  task :stop do 
      kill_processes_matching "go-admin"
  end
  task :restart do
  	restart_cmd = "#{current_release}/go-admin"
  	run "cd #{current_release} && nohup sh -c '#{restart_cmd} -http=127.0.0.1:8082 &' > #{application}-nohup.out"
  end
end

def kill_processes_matching(name)
  begin
    run "killall -q #{name}"
  rescue Exception => error
    puts "No processes to kill"
  end
end

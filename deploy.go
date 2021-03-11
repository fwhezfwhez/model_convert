package model_convert

import (
	"fmt"
	"path"
)

type Param struct {
	PreExports string // 预环境变量配置

	Gopath string // gopath

	RepoGitAddr string // 仓库git地址
	BranchName  string // 分支名

	GenerateToPath string // 产出在哪个目录

	Srvname string // 服务名

	Hosts []string // 发布到哪些服务节点
}

// 产出自动化部署的nginx配置
// 包含了节点的nginx代理，健康检查路由
//
// FAQ:
// 1. 如何接入https
// 使用腾讯云自带的http转https，服务器保持80简单
// 2. 如何做负载均衡
// 使用腾讯云域名负载均衡
// 3. 如何做到http服务高可用
// 使用腾讯云健康检查，检查路径为： https://domain/srvName.status   返回2xx则为服务正常，4xx,5xx则为服务挂起。每次服务重启时，需要做到以下两部：
// - 删除 healthCheckPath/srvName.status文件
// - 部署脚本睡眠2个健康检查周期
// - 服务跑起来后，重新创建srvName.status文件
func GenerateDeployNginxConf(domain string, appPort string, appName string, logabsfile string, healthCheckPath string) string {
	domain = GetDefault(domain, "www.baidu.com")
	appPort = GetDefault(appPort, ":8080")
	appName = GetDefault(appName, "hero")
	logabsfile = GetDefault(logabsfile, fmt.Sprintf("/data/log/nginx/%s.log", domain))
	healthCheckPath = GetDefault(healthCheckPath, "/data/nginx-alives")

	var tmpl = `
server {
    listen      80;
    server_name ${domain};
    error_log ${logabsfile};

    # request header
    proxy_read_timeout 3200;
    proxy_send_timeout 3200;
    proxy_set_header   Host             $http_host;
    proxy_set_header   Cookie           $http_cookie;
    proxy_set_header   X-Real-IP        $remote_addr;
    proxy_set_header   X-Forwarded-Proto    $scheme;
    proxy_set_header   X-Forwarded-For  $proxy_add_x_forwarded_for;

    location / {
        proxy_pass http://localhost${appPort};
    }

    location /${appName}.status {
        root ${healthCheckPath};
    }
    error_page 404 /404.html;
        location = /40x.html {
    }
    error_page 500 502 503 504 /50x.html;
        location = /50x.html {
    }

}
`

	rs := replaceAll(tmpl, map[string]interface{}{
		"${domain}":          domain,
		"${appPort}":         appPort,
		"${appName}":         appName,
		"${logabsfile}":      logabsfile,
		"${healthCheckPath}": healthCheckPath,
	})
	return rs
}

func GenerateDeploySupervisorIni(appName string, workDir string, commandArgs string, logPath string) string {
	appName = GetDefault(appName, "hero")
	workDir = GetDefault(workDir, "/home/web/projects/xyx_srv/")
	command := path.Join(workDir, appName)
	commandArgs = GetDefault(commandArgs, "-mode 'pro'")
	logPath = GetDefault(logPath, fmt.Sprintf("/data/log/%s/%s.log", appName, appName))

	tmpl := `
[program:${appName}]
directory=${workDir}
command=${command} ${commandArgs}
stdout_logfile=${logPath}
stdout_logfile_backups=5
redirect_stderr=true
autostart=true
autorestart=true
`

	rs := replaceAll(tmpl, map[string]interface{}{
		"${appName}":     appName,
		"${workDir}":     workDir,
		"${command}":     command,
		"${commandArgs}": commandArgs,
		"${logPath}":     logPath,
	})

	return rs
}

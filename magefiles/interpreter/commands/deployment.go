package commands

import (
	"os"
	"path"
	"text/template"

	"github.com/magefile/mage/sh"

	"github.com/vindex10/devex-thing/magefiles/common"
)

type DeploymentInitArgs struct {
	ImageName     string
	ImageVersion  string
	AppName       string
	Replicas      int
	ContainerPort int
	ImageSource   string
	ImageRegistry string
}

func NewDeploymentInitArgs(ImageName string, ImageVersion string) DeploymentInitArgs {
	return DeploymentInitArgs{
		ImageName:     ImageName,
		ImageVersion:  ImageVersion,
		AppName:       ImageName,
		Replicas:      1,
		ContainerPort: 8080,
		ImageSource:   common.IMAGE_SOURCE + "/" + ImageName,
		ImageRegistry: common.ARTIFACT_REGISTRY,
	}
}

type DeploymentInit Command[DeploymentInitArgs]

func (c DeploymentInit) Apply(deployment string, args DeploymentInitArgs) error {
	writeDeploymentTpl(deployment, args)
	return nil
}

func writeDeploymentTpl(deployment string, args DeploymentInitArgs) {
	tpl, _ := template.New("tpl-deployment").Parse(`
apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{ .Deployment }}
  labels:
    app: {{ .P.AppName }}
spec:
  replicas: {{ .P.Replicas }}
  selector:
    matchLabels:
      app: {{ .P.AppName }}
  template:
    metadata:
      labels:
        app: {{ .P.AppName }}
    spec:
      containers:
      - name: {{ .P.AppName }}
        image: {{ .ImageFullName }}
        ports:
        - containerPort: {{ .P.ContainerPort }}
`[1:])
	imageFullName := common.AssembleImageFullName(args.ImageRegistry, args.ImageName, args.ImageVersion)
	data := struct {
		Deployment    string
		ImageFullName string
		P             DeploymentInitArgs
	}{
		Deployment:    deployment,
		ImageFullName: imageFullName,
		P:             args,
	}
	deploymentDir := path.Join(common.DEPLOYMENTS_DIR, deployment)
	os.MkdirAll(deploymentDir, 0755)
	fout, _ := os.Create(path.Join(deploymentDir, common.DEPLOYMENT_FILE))
	defer fout.Close()
	tpl.Execute(fout, data)
}

func (c DeploymentInit) ApplyFromStr(deployment string, args string) error {
	return CommandApplyFromStr(c, deployment, args)
}

type DeploymentDelete Command[DummyArgs]

func (c DeploymentDelete) Apply(deployment string, args DummyArgs) error {
	sh.Rm(path.Join(common.DEPLOYMENTS_DIR, deployment))
	return nil
}

func (c DeploymentDelete) ApplyFromStr(deployment string, args string) error {
	return CommandApplyFromStr(c, deployment, args)
}

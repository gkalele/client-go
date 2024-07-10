.PHONY: all dockerbuild pushimage deploy_image push deploy_binary image

all: dockerbuild flannel-probe

TMPFOLDER:=/tmp/ramdisk
TARFILE=probe.tar
VERSION:=`cat version`

flannel-probe: *.go
	go build -v -o flannel-probe .
	strip flannel-probe

push: flannel-probe
	rsync -av flannel-probe orch-esx-3052.tetrationanalytics.com:/tmp/flannel-probe

deploy_binary: push
	rsync -av deploybinary.sh orch-esx-3052.tetrationanalytics.com:/tmp/deploybinary.sh
	ssh orch-esx-3052.tetrationanalytics.com "sudo -u tetinstall -i /tmp/deploybinary.sh"
	kubectl rollout restart ds/flannel-probe

dockerbuild:
	docker build -t flannel-probe:1.$(VERSION) docker/

image: flannel-probe dockerbuild
	docker save flannel-probe:1.$(VERSION) > $(TMPFOLDER)/$(TARFILE)
	ls -lh $(TMPFOLDER)/$(TARFILE)

pushimage: image
	rsync -avz $(TMPFOLDER)/$(TARFILE) orch-esx-3052.tetrationanalytics.com:/tmp/$(TARFILE)

deploy_image: pushimage
	scp deploy.sh orch-esx-3052.tetrationanalytics.com:/tmp/deploy.sh
	ssh orch-esx-3052.tetrationanalytics.com "sudo -u tetinstall -i /tmp/deploy.sh"

increment:
	@expr `cat version` + 1 > version
	@cat version

daemonset:
	jinja2 -D version=`cat version` probe-daemonset.yml.j2 > probe-daemonset.yml
	kubectl delete ds flannel-probe || true
	kubectl create -f probe-daemonset.yml

# noos-demo make utility
# derived from https://github.com/usbarmory/go-boot/blob/main/Makefile

BUILD_TAGS = linkcpuinit,linkramsize,linkramstart,linkprintk
SHELL = /bin/bash
APP ?= noosdemo

IMAGE_BASE := 10000000
TEXT_START := $(shell echo $$((16#$(IMAGE_BASE) + 16#10000)))
GOFLAGS := -tags ${BUILD_TAGS} -trimpath -ldflags "-s -w -E cpuinit -T $(TEXT_START) -R 0x1000"

OVMFCODE ?= OVMF_CODE.fd
OVMFVARS ?= OVMF_VARS.fd

QEMU ?= qemu-system-x86_64 \
        -enable-kvm -cpu host,invtsc=on -m 8G \
        -drive file=fat:rw:$(CURDIR) \
        -drive if=pflash,format=raw,readonly,file=$(OVMFCODE) \
        -drive if=pflash,format=raw,file=$(OVMFVARS) \
        -global isa-debugcon.iobase=0x402 \
        -serial stdio -display sdl

.PHONY: clean

qemu: $(APP).efi
	@if [ "${QEMU}" == "" ]; then \
		echo 'qemu not available for this target'; \
		exit 1; \
	fi

	cp ${APP}.efi efi/boot/bootx64.efi
	$(QEMU)

check_tamago:
	@if [ "${TAMAGO}" == "" ] || [ ! -f "${TAMAGO}" ]; then \
		echo 'You need to set the TAMAGO variable to a compiled version of https://github.com/usbarmory/tamago-go'; \
		exit 1; \
	fi

clean:
	@rm -fr $(APP)-tamago $(APP).efi $(APP)-tinygo

tamago: check_tamago
	GOOS=tamago GOARCH=amd64 $(TAMAGO) build $(GOFLAGS) -o ${APP}-tamago

$(APP).efi: tamago
	objcopy \
		--strip-debug \
		--target efi-app-x86_64 \
		--subsystem=efi-app \
		--image-base 0x$(IMAGE_BASE) \
		--stack=0x10000 \
		${APP}-tamago ${APP}.efi
	printf '\x26\x02' | dd of=${APP}.efi bs=1 seek=150 count=2 conv=notrunc,fsync # adjust Characteristics

	rm ${APP}-tamago

check_tinygo:
	@if ! command -v tinygo &> /dev/null ; then \
		echo 'You need to install the tinygo compiler'; \
		exit 1; \
	fi

tinygo: check_tinygo
	tinygo build -tags noos,noasm -o ${APP}-tinygo

check_tinydisplay:
	@if ! command -v tinydisplay &> /dev/null ; then \
		echo 'You need to install the tinydisplay display simulator from https://github.com/sago35/tinydisplay'; \
		exit 1; \
	fi

tinydisplay: check_tinydisplay
	tinydisplay &
	go run -tags tinygo,noos,noasm .

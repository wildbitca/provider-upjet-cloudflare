FROM alpine:3.23.2
RUN apk --no-cache add ca-certificates bash

ARG TARGETOS
ARG TARGETARCH

ADD "bin/${TARGETOS}_${TARGETARCH}/provider" /usr/local/bin/provider

ENV USER_ID=65532

# Setup OpenTofu + provider environment

## Provider-dependent configuration
ARG OPENTOFU_VERSION
ARG TERRAFORM_PROVIDER_SOURCE
ARG TERRAFORM_PROVIDER_VERSION
ARG TERRAFORM_PROVIDER_DOWNLOAD_NAME
ARG TERRAFORM_NATIVE_PROVIDER_BINARY
ARG TERRAFORM_PROVIDER_DOWNLOAD_URL_PREFIX

## End of - Provider-dependent configuration

ENV PLUGIN_DIR=/terraform/provider-mirror/registry.terraform.io/${TERRAFORM_PROVIDER_SOURCE}/${TERRAFORM_PROVIDER_VERSION}/${TARGETOS}_${TARGETARCH}
ENV TF_CLI_CONFIG_FILE=/terraform/.terraformrc
ENV TF_FORK=0

RUN mkdir -p ${PLUGIN_DIR}

ADD https://github.com/opentofu/opentofu/releases/download/v${OPENTOFU_VERSION}/tofu_${OPENTOFU_VERSION}_${TARGETOS}_${TARGETARCH}.zip /tmp
ADD ${TERRAFORM_PROVIDER_DOWNLOAD_URL_PREFIX}/${TERRAFORM_PROVIDER_DOWNLOAD_NAME}_${TERRAFORM_PROVIDER_VERSION}_${TARGETOS}_${TARGETARCH}.zip /tmp
ADD terraformrc.hcl ${TF_CLI_CONFIG_FILE}

RUN unzip /tmp/tofu_${OPENTOFU_VERSION}_${TARGETOS}_${TARGETARCH}.zip -d /usr/local/bin \
  && chmod +x /usr/local/bin/tofu \
  && ln -s /usr/local/bin/tofu /usr/local/bin/terraform \
  && rm /tmp/tofu_${OPENTOFU_VERSION}_${TARGETOS}_${TARGETARCH}.zip \
  && unzip /tmp/${TERRAFORM_PROVIDER_DOWNLOAD_NAME}_${TERRAFORM_PROVIDER_VERSION}_${TARGETOS}_${TARGETARCH}.zip -d ${PLUGIN_DIR} \
  && chmod +x ${PLUGIN_DIR}/* \
  && rm /tmp/${TERRAFORM_PROVIDER_DOWNLOAD_NAME}_${TERRAFORM_PROVIDER_VERSION}_${TARGETOS}_${TARGETARCH}.zip \
  && chown -R ${USER_ID}:${USER_ID} /terraform
# End of - Setup OpenTofu + provider environment

# Provider controller needs these environment variables at runtime
ENV TERRAFORM_VERSION=${OPENTOFU_VERSION}
ENV TERRAFORM_PROVIDER_SOURCE=${TERRAFORM_PROVIDER_SOURCE}
ENV TERRAFORM_PROVIDER_VERSION=${TERRAFORM_PROVIDER_VERSION}
ENV TERRAFORM_NATIVE_PROVIDER_PATH=${PLUGIN_DIR}/${TERRAFORM_NATIVE_PROVIDER_BINARY}

USER ${USER_ID}
EXPOSE 8080

ENTRYPOINT ["provider"]

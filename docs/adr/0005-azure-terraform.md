# ADR 0005 - Deploy on Azure with Terraform

* Status: **Accepted**
* Date: 2026-11-12
* Deciders: the whole team

## Context

AREA has to be deployable by someone who was not in the room. Locally that stays Docker Compose. The hosted environment should be the same idea: a small Azure setup, written in Terraform, applied from GitHub.

The team is small. Planning and reviews already live on GitHub. The hosted app has to run the NestJS API, PostgreSQL, and the web build. It does not have to survive a region failure or a traffic spike.

## Decision

Host the deployed app on **Azure**, and manage that hosting with **Terraform**.

The setup stays small:

* One resource group.
* **Azure Container Apps** for the NestJS API. One container, and it can scale to zero when nobody is calling it.
* **Azure Database for PostgreSQL Flexible Server**, burstable, one server. Same Postgres as ADR 0003.
* **Azure Container Registry** for the API image.
* The web build (Wasm) is static files on **Azure Storage**. The Android app is an installable artifact from CI, not a hosted service.
* Terraform state in an Azure Storage account, so an apply is locked while another one is running.
* GitHub Actions, logged into Azure with a federated credential (no cloud password stored in the repo):
  * on a pull request: lint, test, build, and `terraform plan`
  * on main: build the image, `terraform apply`, deploy

A Redis cache for BullMQ is not part of this setup. We add it when a worker actually needs it.

Kubernetes, Front Door, a second region and private networking are out. Any of those is a new ADR.

### Alternatives rejected

* **Clicking through the Azure portal.** Fast the first time. Nobody can rebuild it afterwards, and the subject asks for a reproducible deploy.
* **AKS.** A cluster is something we would have to operate. We do not have the traffic, and we do not want the on-call.
* **Bicep.** It is native to Azure. Terraform is the one that gives us `plan` on the pull request and `apply` on main, and the files stay readable if the host ever changes.

## Consequences

### Positive

* An infra change shows up as a plan on the pull request before it runs.
* GitHub stays the one place we look: planning, CI and deploy.
* What runs in Azure is a container plus Postgres, which is what we already run locally. Docker Compose stays the dev setup.
* Container Apps scales to zero, so a quiet week is not a cluster bill.

### Negative

1. **Postgres does not scale to zero.** The burstable server keeps billing while the project sits still. We accept that and stay on the smallest SKU.
2. **One region, one server.** A bad apply or an Azure incident takes the demo down. We do not add a second region in advance. The Terraform state lives in the storage account, and that account is the backup of the setup.
3. **Two tools to learn.** Terraform state locks and Container Apps revisions both bite. The module list stays short so the pipeline is something the whole team can read.

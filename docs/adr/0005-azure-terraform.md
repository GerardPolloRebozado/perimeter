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

## Cost

List prices for **France Central**, from the Azure retail prices API on 13 Nov 2026. A month is 730 hours. The API is left running and Container Apps is scaled to zero.

| Piece | Meter | Month |
|---|---|---|
| Postgres Flexible Server, Burstable B1MS | $0.019 per hour | $13.87 |
| Postgres disk, 32 GB | $0.133 per GB | $4.26 |
| Container Registry, Basic | $0.1666 per day | $5.07 |
| Blob storage for the Wasm site and Terraform state | $0.0192 per GB, Hot LRS | a few cents |
| Container Apps, inside the free grant | $0 | $0 |
| **Total** | | **about $23** |

The $100 Azure for Students credit covers that for about four months. The same offer, for a new Azure customer and for 12 months, includes 750 hours of this B1MS server with 32 GB of storage and 32 GB of backup, plus one Standard registry. Those two lines then drop to $0 and the credit stays put. Container Apps is always free for the first 180,000 vCPU-seconds, 360,000 GB-seconds and 2 million requests each month.

One replica left running all month is about $10 if Azure counts it as idle, and about $34 if the Node process stays above the idle threshold (0.01 vCPU). That second case spends the $100 in under two months. High availability, a second region and Redis are not in this table.

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

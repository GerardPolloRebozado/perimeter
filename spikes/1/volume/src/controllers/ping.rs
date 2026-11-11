#![allow(clippy::missing_errors_doc)]
#![allow(clippy::unnecessary_struct_initialization)]
#![allow(clippy::unused_async)]
use loco_rs::prelude::*;
use serde::Serialize;
use serde_json::json;
use utoipa::ToSchema;

#[derive(Serialize, Debug, ToSchema)]
pub struct Ping {
    pub status: String,
}

/// Ping endpoint
///
/// This endpoint will return status ok if the api is working
#[utoipa::path(
    get,
    path = "/ping",
    responses(
        (status = 200, description = "This endpoint will return status ok if the api is working", body = Ping),
    ),
)]
#[debug_handler]
pub async fn index(State(_ctx): State<AppContext>) -> Result<Response> {
    format::json(json!(Ping {
        status: String::from("ok")
    }))
}

pub fn routes() -> Routes {
    Routes::new().prefix("ping/").add("/", get(index))
}

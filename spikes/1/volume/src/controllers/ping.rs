#![allow(clippy::missing_errors_doc)]
#![allow(clippy::unnecessary_struct_initialization)]
#![allow(clippy::unused_async)]
use loco_rs::prelude::*;
use serde_json::json;

#[debug_handler]
pub async fn index(State(_ctx): State<AppContext>) -> Result<Response> {
format::json(json!({
        "status": "ok"
    }))
}

pub fn routes() -> Routes {
    Routes::new()
        .prefix("ping/")
        .add("/", get(index))
}

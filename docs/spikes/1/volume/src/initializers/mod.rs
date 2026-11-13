use async_trait::async_trait;
use axum::Router;
use loco_rs::{
    app::{AppContext, Initializer},
    Result,
};
use utoipa::OpenApi;
use utoipa_swagger_ui::SwaggerUi;

#[derive(OpenApi)]
#[openapi(
    paths(
        crate::controllers::ping::index,
        crate::controllers::auth::register,
        crate::controllers::auth::login,
        crate::controllers::auth::current,
    ),
    components(schemas(
        crate::controllers::ping::Ping,
        crate::views::auth::LoginResponse,
        crate::views::auth::CurrentResponse,
        crate::models::users::RegisterParams,
        crate::models::users::LoginParams,
    ))
)]
struct ApiDoc;

/// Loco initializer for OpenAPI docs, served with Swagger UI at `/api-docs`.
pub struct OpenapiInitializer;

#[async_trait]
impl Initializer for OpenapiInitializer {
    fn name(&self) -> String {
        "openapi".to_string()
    }

    async fn after_routes(&self, router: Router, _ctx: &AppContext) -> Result<Router> {
        let openapi = ApiDoc::openapi();
        let swagger_ui: Router = SwaggerUi::new("/api-docs")
            .url("/api-docs/openapi.json", openapi)
            .into();
        Ok(router.merge(swagger_ui))
    }
}

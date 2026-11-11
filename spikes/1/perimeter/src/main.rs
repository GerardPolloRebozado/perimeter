use dioxus::prelude::*;

#[derive(Debug, Clone, Routable, PartialEq)]
#[rustfmt::skip]
enum Route {
    #[layout(Navbar)]
    #[route("/")]
    Home {},
}

const FAVICON: Asset = asset!("/assets/favicon.ico");
const MAIN_CSS: Asset = asset!("/assets/main.css");
const HEADER_SVG: Asset = asset!("/assets/header.svg");
const TAILWIND_CSS: Asset = asset!("/assets/tailwind.css");

fn main() {
    dioxus::launch(App);
}

#[component]
fn App() -> Element {
    rsx! {
        document::Link { rel: "icon", href: FAVICON }
        document::Link { rel: "stylesheet", href: MAIN_CSS } document::Link { rel: "stylesheet", href: TAILWIND_CSS }
        Router::<Route> {}
    }
}

/// Home page
#[component]
fn Home() -> Element {
    rsx! {
        PingComponent {}
    }
}

/// Shared navbar component.
#[component]
fn Navbar() -> Element {
    rsx! {
        div {
            id: "navbar",
            Link {
                to: Route::Home {},
                "Home"
            }
        }

        Outlet::<Route> {}
    }
}

#[component]
fn PingComponent() -> Element {
    let response = use_signal(|| String::new());

    // This runs on the server during SSR.
    // The `?` operator triggers Dioxus Suspense and waits for it to resolve
    // before sending the completed HTML down to the browser.
    let ping_data = use_loader(move || async move {
        let text = reqwest::get("http://localhost:5150/ping")
            .await?
            .text()
            .await?;
        Ok::<String, reqwest::Error>(text)
    })?;

    rsx! {
        div {
            id: "echo",
            h4 { "Server-Side Ping Echo" }

            p {
                "Ping response from server: "
                i { "{ping_data}" }
            }

            if !response().is_empty() {
                p {
                    "You typed: "
                    i { "{response}" }
                }
            }
        }
    }
}

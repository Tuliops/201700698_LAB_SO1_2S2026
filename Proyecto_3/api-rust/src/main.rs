use axum::{routing::post, Json, Router};
use serde::{Deserialize, Serialize};
use std::net::SocketAddr;

#[derive(Debug, Deserialize, Serialize, Clone)]
struct ReportPayload {
    country: String,
    warplanes_in_air: i32,
    warships_in_water: i32,
    timestamp: String,
}

async fn handle_report(Json(payload): Json<ReportPayload>) -> String {
    println!("[Rust] Recibido reporte: {:?}", payload);

    let client = reqwest::Client::new();
    let go_url = "http://go-d1-service:8081/recibir";

    match client.post(go_url).json(&payload).send().await {
        Ok(res) => {
            println!("[Rust] Go D1 respondió con status: {}", res.status());
            "Reporte procesado y enviado a Go D1".to_string()
        }
        Err(e) => {
            println!("[Rust] Error llamando a Go D1: {}", e);
            format!("Error llamando a Go D1: {}", e)
        }
    }
}

#[tokio::main]
async fn main() {
    let app = Router::new()
        .route("/grpc-201700698", post(handle_report));

    let addr = SocketAddr::from(([0, 0, 0, 0], 8080));
    println!("[Rust] API REST escuchando en {}", addr);
    
    let listener = tokio::net::TcpListener::bind(addr).await.unwrap();
    axum::serve(listener, app).await.unwrap();
}


//manejo de peticiones https
use actix_web::{web,App, HttpResponse, HttpServer, Responder};
// serialializacion y Deserializacion de datos 
use serde::{Deserialize, Serialize};



// gRPC generado por libreira tonic a partir del .protc

pub mod wartweets{
    tonic::include_proto!("wartweets");
}
use wartweets::war_report_service_client::WarReportServiceClient;
use wartweets:: WarReportRequest ;

// Definimos la estructura del reporte militar
#[derive(Serialize,Deserialize,Debug)]
struct WarReport {
    country:String,
    warplanes_in_air:i32,
    warships_in_water: i32,
    timestamp:String,
}
// funcion convertir a string enum del .protc

fn parse_country(country_str: &str) -> i32{
    match country_str.to_uppercase().as_str(){
        "USA" => 1,
        "RUS" => 2,
        "CHN" => 3,
        "ESP" => 4,
        "GTM" => 5,
        _ => 0,
        
    }
}

//Funcion que  resive los datos y los procesa
async fn receive_report(report: web::Json<WarReport>) -> impl Responder {
    println!(
        "[Rust REST] Reporte recibido de {}! Aviones: {}, Barcos: {}, Hora: {}",
        report.country, report.warplanes_in_air, report.warships_in_water, report.timestamp
    );

    // Conexion al servidor gRPC de Go corriendo en el puerto 50051
    let grpc_addr = "http://127.0.0.1:50051";
    
    //  enviar el reporte mediante gRPC
    match WarReportServiceClient::connect(grpc_addr).await {
        Ok(mut client) => {
            let request = tonic::Request::new(WarReportRequest {
                country: parse_country(&report.country),
                warplanes_in_air: report.warplanes_in_air,
                warships_in_water: report.warships_in_water,
                timestamp: report.timestamp.clone(),
            });

            match client.send_report(request).await {
                Ok(response) => {
                    println!(" gRPC enviado con exito a Go. Respuesta: {:?}", response.get_ref().status);
                    HttpResponse::Ok().json(serde_json::json!({
                        "status": "success",
                        "message": "Reporte procesado en Rust y reenviado a Go via gRPC"
                    }))
                }
                Err(e) => {
                    eprintln!(" Error al llamar el servicio gRPC: {}", e);
                    HttpResponse::InternalServerError().json(serde_json::json!({
                        "status": "error",
                        "message": "Falla al comunicarse con el servidor gRPC en Go"
                    }))
                }
            }
        }
        Err(e) => {
            eprintln!(" No se pudo conectar al servidor gRPC en Go ({grpc_addr}): {}", e);
            HttpResponse::InternalServerError().json(serde_json::json!({
                "status": "error",
                "message": "No hay conexion con el backend gRPC"
            }))
        }
    }
}


#[actix_web::main]
async fn main() -> std::io::Result<()> {
    let host = "127.0.0.1";
    let port = 8080;

    println!("Servidor REST en Rust corriendo en http://{}:{}", host, port);

    HttpServer::new(|| {
        App::new()
            .route("/report", web::post().to(receive_report))
    })
    .bind((host, port))?
    .run()
    .await
}
use std::io::{Read, Write};
use std::net::TcpStream;

// GPU monitoring
use nvml_wrapper::enum_wrappers::device::TemperatureSensor;
use nvml_wrapper::error::NvmlError as NvmlError;
use nvml_wrapper::NVML;
// System monitoring
use sysinfo::{Disks, System,};
// Timer
use std::{thread, time};
// json
use serde_json::json;
use serde_json::Value;

fn main() {
    
    loop {
        let mut stream = TcpStream::connect("172.31.224.143:10772").unwrap();
        println!("Connected to the server");
        let result: Result<serde_json::Value, std::string::String> = gather_metrics();
        let result_string = match result {
            Ok(value) => value.to_string(),
            Err(e) => e,
        };
       
        match stream.write_all(result_string.as_bytes()) {
            Ok(_) => println!("Data sent successfully"),
            Err(e) => {
                eprintln!("Failed to write to socket: {}", e);
                return;
            }
        }

        println!("Sent message: {}", result_string);

        // Now read the response in a similar way to how we refined the server
        // let mut buffer = Vec::new();
        // let mut temp_buf = [0; 1024]; // Temporary buffer to read chunks

        // loop {
        //     match stream.read(&mut temp_buf) {
        //         Ok(0) => {
        //             // No more data to read
        //             break;
        //         },
        //         Ok(n) => {
        //             buffer.extend_from_slice(&temp_buf[..n]);

        //             // If less than 1024 bytes are read, it might be the end of the message
        //             if n < 1024 {
        //                 break;
        //             }
        //         },
        //         Err(e) => {
        //             // No more data available for now
        //             break;
        //         }
        //     }
        // }

        // let received = String::from_utf8_lossy(&buffer);
        // println!("\nReceived echo: {}", received);
        
        thread::sleep(time::Duration::from_secs(2));
        // thread::sleep(time::Duration::from_secs(15 * 60)); // Sleep for 15 minutes
    }
}

pub fn gather_metrics() -> Result<serde_json::Value, String> {
    let mut system = System::new_all();
    system.refresh_all();

    let memory_used = system.used_memory();
    let total_memory = system.total_memory();

    let disks = Disks::new_with_refreshed_list();
    // Compute disk_total and disk_available and disk_used from disks
    let disk_total = disks.iter().map(|disk| disk.total_space()).sum::<u64>();
    let disk_available = disks.iter().map(|disk| disk.available_space()).sum::<u64>();
    let disk_used = disk_total - disk_available;

    let cpu_usage = system.cpus().iter().map(|cpu| cpu.cpu_usage()).sum::<f32>() / system.cpus().len() as f32;


    // Get GPU information
    let nvml = NVML::init().unwrap();
    let gpu_info = match dump_all_gpu_info(&nvml) {
        Ok(info) => info,
        Err(e) => return Err(format!("Failed to read GPU info: {}", e)),
    };

    let metrics = json!({
        "memory": {
            "used": memory_used,
            "total": total_memory
        },
        "disk": {
            "available": disk_available,
            "used": disk_used,
            "total": disk_total
        },
        "gpu": gpu_info,
        "cpu_usage": cpu_usage,
    });

    Ok(metrics)
}

pub fn read_gpu_info(device: &nvml_wrapper::Device) -> Result<Value, NvmlError> {
    let name = device.name()?;
    let id = device.index()?;
    let compute_capability = device.cuda_compute_capability()?;
    let utilization_rates = device.utilization_rates()?;
    let memory_info = device.memory_info()?;
    let temperature = device.temperature(TemperatureSensor::Gpu)?;

    let gpu_info = serde_json::json!({
        "id": id,
        "name": name,
        "compute_capability": format!("{}.{}", compute_capability.major, compute_capability.minor),
        "utilization_gpu": utilization_rates.gpu,
        "memory_used_mb": memory_info.used / 1024 / 1024,
        "memory_total_mb": memory_info.total / 1024 / 1024,
        "temperature_c": temperature
    });

    Ok(gpu_info)
}

pub fn dump_all_gpu_info(nvml: &nvml_wrapper::NVML) -> Result<Value, NvmlError> {
    let device_count = nvml.device_count()?;
    let mut gpu_infos = Vec::new();

    for i in 0..device_count {
        let device = nvml.device_by_index(i)?;
        let gpu_info = read_gpu_info(&device)?;
        gpu_infos.push(gpu_info);
    }

    Ok(serde_json::Value::Array(gpu_infos))
}

// pub fn read_gpu_info(device: &nvml_wrapper::Device) -> Result<String, NvmlError> {
//     let name = device.name()?;
//     let id = device.index()?;
//     let compute_capability = device.cuda_compute_capability()?;
//     let utilization_rates = device.utilization_rates()?;
//     let memory_info = device.memory_info()?;
//     let temperature = device.temperature(TemperatureSensor::Gpu)?;

//     Ok(format!(
//         "[{}]: {} | {}.{} | {} % | {} / {} MB | {}°C |",
//         id,
//         name,
//         compute_capability.major,
//         compute_capability.minor,
//         utilization_rates.gpu,
//         memory_info.used / 1024 / 1024,
//         memory_info.total / 1024 / 1024,
//         temperature
//     ))
// }

// pub fn dump_all_gpu_info(nvml: &nvml_wrapper::NVML) -> Result<String, NvmlError> {
//     let device_count = nvml.device_count()?;
//     let mut all_gpu_info = String::new();

//     for i in 0..device_count {
//         let device = nvml.device_by_index(i)?;
//         let gpu_info = read_gpu_info(&device)?;
//         all_gpu_info.push_str(&gpu_info);
//         all_gpu_info.push_str("\n"); // Adds a newline for each GPU entry
//     }

//     Ok(all_gpu_info)
// }
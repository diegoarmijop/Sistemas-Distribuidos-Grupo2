# monitor.py
import subprocess
import time
import requests
import sys

def watch_logs():
    subprocess.run(["docker-compose", "logs", "-f"])

def watch_metrics():
    while True:
        print("=== Métricas del Sistema ===")
        subprocess.run(["docker", "stats", "--no-stream"])
        time.sleep(5)

def check_health():
    print("=== Estado de los Servicios ===")
    subprocess.run(["docker-compose", "ps"])
    
    print("\n=== Estado de RabbitMQ ===")
    try:
        response = requests.get("http://localhost:15672/api/health/checks/alarms", 
                              auth=('guest', 'guest'))
        print(response.json())
    except Exception as e:
        print(f"Error checking RabbitMQ: {e}")
    
    print("\n=== Estado de las Bases de Datos ===")
    for port in [5432, 5434, 5435, 5436]:
        subprocess.run(["docker", "exec", "postgres_container", 
                       "pg_isready", "-h", "localhost", "-p", str(port)])

def main():
    while True:
        print("\n=== Monitor del Sistema ===")
        print("1. Ver logs en tiempo real")
        print("2. Monitorear métricas")
        print("3. Verificar estado de servicios")
        print("4. Salir")
        
        choice = input("Seleccione una opción: ")
        
        if choice == "1":
            watch_logs()
        elif choice == "2":
            watch_metrics()
        elif choice == "3":
            check_health()
        elif choice == "4":
            sys.exit(0)
        else:
            print("Opción inválida")

if __name__ == "__main__":
    main()
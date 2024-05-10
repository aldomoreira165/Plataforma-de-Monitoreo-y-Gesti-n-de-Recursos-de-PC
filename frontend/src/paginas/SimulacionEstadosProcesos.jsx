import { React, useState } from 'react'
import Swal from 'sweetalert2';
import { Graphviz } from 'graphviz-react';
import Button from 'react-bootstrap/Button';
import '../assets/css/main.css'

export default function SimulacionEstadosProcesos() {
    const [pid, setPid] = useState("");
    const [historial, setHistorial] = useState([]);
    const [estados, setEstados] = useState([
        { stateName: "new", activo: false },
        { stateName: "ready", activo: false },
        { stateName: "running", activo: false }
    ]);

    const [edges, setEdges] = useState(["{new -> ready;}", "{ready -> running;}"]);

    const inputId = document.querySelector('#input-id');

    const mensajeExito = (mensaje) => {
        Swal.fire({
            position: "center",
            icon: "success",
            title: mensaje,
            showConfirmButton: false,
            timer: 1500
        });
    }

    const mensajeError = (mensaje) => {
        Swal.fire({
            position: "center",
            icon: "error",
            title: mensaje,
            showConfirmButton: false,
            timer: 1500
        });
    }

    const enviarEstado = (estado) => {
        fetch('/api/insertar_estado', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                id_proceso: parseInt(inputId.value),
                estado: estado,
            }),
        })
        .then(response => {
            if (!response.ok) {
                throw new Error('Error al enviar el estado al servidor');
            }
        })
        .catch(error => {
            console.error(error);
            mensajeError('Error al enviar el estado al servidor');
        });
    };

    const obtenerHistorialProceso = () => {
        fetch(`/api/obtener_historial_estado?pid=${inputId.value}`)
            .then(response => response.json())
            .then(data => {
                setHistorial(data);
                console.log(historial);
            })
            .catch(error => {
                console.error(error);
                mensajeError('Error al obtener el historial del proceso');
            });
    };

    const iniciarProceso = () => {
        fetch('/api/start_process')
            .then(response => response.text())
            .then(data => {
                setPid(data)
                inputId.value = pid;
                actualizarEstado("running");
                enviarEstado("new");
                enviarEstado("ready");
                enviarEstado("running");
                mensajeExito(`Proceso con id ${pid} creado correctamente`);

                setEdges([
                    "{new -> ready;}",
                    "{ready -> running;}",
                ]);
            })
            .catch(error => {
                mensajeError('Error al crear el proceso');
                console.log(error)
            });
    }

    const matarProceso = () => {
        fetch(`/api/kill_process?pid=${inputId.value}`, {
            method: 'DELETE',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ pid: inputId.value })
        })
            .then(response => response.text())
            .then(data => {
                setEstados([
                    { stateName: "new", activo: false },
                    { stateName: "ready", activo: false },
                    { stateName: "running", activo: false }
                ]);
                enviarEstado("terminated");
                setEdges([
                    "{new -> ready;}",
                    "{ready -> running;}",
                ]);
                mensajeExito(`Proceso con id ${inputId.value} eliminado correctamente`);
            })
            .catch(error => {
                mensajeError('Error al eliminar el proceso');
                console.log(error)
            });
    }

    const resumirProceso = () => {
        fetch(`/api/resume_process?pid=${inputId.value}`)
            .then(response => response.text())
            .then(data => {
                mensajeExito(`Proceso con id ${inputId.value} continuado correctamente`);
                actualizarEstado("running");
                enviarEstado("resumed");
                setEdges([
                    "{new -> ready;}",
                    "{ready -> running [color=\"red\"];}",
                    "{running -> ready;}"
                ]);
            })
            .catch(error => {
                mensajeError('Error al continuar el proceso');
                console.log(error)
            });
    }

    const detenerProceso = () => {
        fetch(`/api/stop_process?pid=${inputId.value}`)
            .then(response => response.text())
            .then(data => {
                mensajeExito(`Proceso con id ${inputId.value} detenido correctamente`);
                actualizarEstado("ready");
                enviarEstado("waiting");
                setEdges([
                    "{new -> ready;}",
                    "{ready -> running;}",
                    "{running -> ready [color=\"red\"];}"
                ]);
            })
            .catch(error => {
                mensajeError('Error al detener el proceso');
                console.log(error)
            });
    }

    const actualizarEstado = (estadoActivo) => {
        setEstados(prevEstados => {
            return prevEstados.map(estado => ({
                ...estado,
                activo: estado.stateName === estadoActivo
            }));
        });
    }

    const graficar = () => {
        const nodes = estados.map(estado => {
            const color = estado.activo ? "yellowgreen" : "deepskyblue";
            return `{${estado.stateName} [label="${estado.stateName}", shape="circle", style="filled", fillcolor="${color}", width=1, height=1];}`;
        });

        const rankdir = "LR";
    
        // Concatenar todos los nodos, aristas y la configuración de dirección
        return `digraph { rankdir=${rankdir}; ${nodes.join(" ")} ${edges.join(" ")} }`;
    };

    const generarDiagramaHistorial = () => {
        const historialNodes = historial.map((estado, index) => {
            return `{${estado.estado}${index} [label="${estado.estado}", shape="circle", style="filled", fillcolor="lightgray", width=1, height=1];}`;
        });
    
        const historialEdges = historial.map((estado, index) => {
            if (index === 0) {
                return ""; // No hay arista para la primera entrada del historial
            } else {
                return `{${historial[index - 1].estado}${index - 1} -> ${estado.estado}${index};}`;
            }
        });
    
        const rankdir = "LR";
    
        const dotString = `digraph { rankdir=${rankdir}; ${historialNodes.join(" ")} ${historialEdges.join(" ")} }`;
    
        return <Graphviz dot={dotString} />;
    };

    return (
        <div className='main-container'>
            <h1 className='main-container__title'>Diagrama de Estados</h1>
            <div className='process-container'>
                <div className='process-container__options'>
                    <div className='process-container__inputbox'>
                        <input id='input-id' className='process-container__input' type="text" />
                    </div>
                    <div className='process-container__buttons'>
                        <Button onClick={iniciarProceso} variant="success" size="lg">New</Button>{' '}
                        <Button onClick={detenerProceso} variant="warning" size="lg">Stop</Button>{' '}
                        <Button onClick={resumirProceso} variant="info" size="lg">Resume</Button>{' '}
                        <Button onClick={matarProceso} variant="danger" size="lg">Kill</Button>{' '}
                        <Button onClick={obtenerHistorialProceso} variant="secondary" size="lg">Historial de Proceso</Button>{' '}
                    </div>
                </div>
                <div className='process-container__graph'>
                    <Graphviz dot={graficar()} />
                </div>
                <div className='process-container__diagram'>
                    {generarDiagramaHistorial()}
                </div>
            </div>
        </div>
    )
}

import React from 'react'
import Container from 'react-bootstrap/Container';
import Nav from 'react-bootstrap/Nav';
import Navbar from 'react-bootstrap/Navbar';
import NavDropdown from 'react-bootstrap/NavDropdown';
import docker from '../assets/img/docker.png'
import '../assets/css/navegacion.css'

export default function Navegacion() {
    return (
        <Navbar collapseOnSelect expand="lg" className="bg-dark" variant="dark">
            <Container>
                <div className='container-name'>
                    <img src={docker} alt="logo-docker"  className='logo'/>
                    <Navbar.Brand href="/"><strong>Plataforma de Monitoreo SO1</strong></Navbar.Brand>
                </div>
                <div className='container-dropdown'>
                    <Navbar.Toggle aria-controls="responsive-navbar-nav" />
                    <Navbar.Collapse id="responsive-navbar-nav">
                        <Nav className="me-auto">
                            <NavDropdown title="Módulos de Kernel" id="collapsible-nav-dropdown">
                                <NavDropdown.Item href="/" className="text-dark" >Monitoreo en tiempo real de memoria RAM y CPU</NavDropdown.Item>
                                <NavDropdown.Item href="/RAMCPUH" className="text-dark">Monitoreo histórico de memoria RAM y CPU</NavDropdown.Item>
                                <NavDropdown.Item href="/ArbolProcesos" className="text-dark">Árbol de Procesos</NavDropdown.Item>
                                <NavDropdown.Item href="/SimulacionEstadosProcesos" className="text-dark">Simulación de Cambio de Estados en los Procesos</NavDropdown.Item>
                            </NavDropdown>
                        </Nav>
                    </Navbar.Collapse>
                </div>
            </Container>
        </Navbar >
    );
}

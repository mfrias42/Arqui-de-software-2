import React, { useState, useEffect } from 'react';
import axios from 'axios';
import '../assets/styles/UploadFile.css';
import { useParams, useNavigate } from 'react-router-dom';

function UploadFile() {
    const params = useParams();
    const courseId = params.courseId;
    const navigate = useNavigate();
    const [file, setFile] = useState(null);
    const [error, setError] = useState(null);

    useEffect(() => {
        console.log('Params:', params);
        console.log('CourseId:', courseId);
        
        if (!courseId) {
            console.error('No se encontró courseId');
            setError('ID de curso no encontrado');
            return;
        }
    }, [courseId]);

    const handleFileChange = (e) => {
        const selectedFile = e.target.files[0];
        setFile(selectedFile);
    };

    const handleSubmit = async (e) => {
        e.preventDefault();

        if (!file) {
            setError('Por favor, selecciona un archivo para subir.');
            return;
        }

        if (!courseId) {
            setError('ID de curso no válido');
            return;
        }

        const reader = new FileReader();
        reader.onloadend = async () => {
            const base64String = reader.result.split(',')[1];

            try {
                console.log('Enviando petición:', {
                    url: `http://localhost:8080/courses/${courseId}/files`,
                    data: {
                        name: file.name,
                        content: base64String.substring(0, 20) + '...',
                        userId: Number(localStorage.getItem('userId'))
                    }
                });

                const response = await axios.post(`http://localhost:8080/courses/${courseId}/files`, {
                    name: file.name,
                    content: base64String,
                    userId: Number(localStorage.getItem('userId'))
                }, {
                    headers: {
                        'Content-Type': 'application/json',
                        Authorization: `Bearer ${localStorage.getItem('token')}`
                    }
                });

                console.log('Respuesta:', response);
                alert('Archivo subido correctamente');
                navigate(`/courses/${courseId}/files`);
            } catch (error) {
                console.error('Error detallado:', error.response || error);
                setError(`Error al subir el archivo: ${error.response?.data?.error || error.message}`);
            }
        };

        reader.readAsDataURL(file);
    };

    if (error) {
        return (
            <div className="upload-container">
                <div className="error-message">{error}</div>
                <button className="back-button" onClick={() => navigate(-1)}>Volver</button>
            </div>
        );
    }

    return (
        <div className="upload-container">
            <button className="back-button" onClick={() => navigate(-1)}>Volver</button>
            <h2>Subir Archivo al Curso {courseId}</h2>
            <form onSubmit={handleSubmit}>
                <div className="file-input-container">
                    <label htmlFor="file-upload" className="custom-file-upload">
                        Seleccionar archivo
                    </label>
                    <input
                        id="file-upload"
                        type="file"
                        onChange={handleFileChange}
                    />
                    {file && <span className="file-name">{file.name}</span>}
                </div>
                <button type="submit" className="upload-button">
                    Subir Archivo
                </button>
            </form>
        </div>
    );
}

export default UploadFile;

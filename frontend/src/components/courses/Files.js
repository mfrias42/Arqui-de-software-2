import React, { useState } from 'react';
import axios from 'axios';
import { useParams } from 'react-router-dom';
import '../assets/styles/FileUploadComponent.css';

const FileUploadComponent = () => {
    const { courseId } = useParams();
    const [selectedFile, setSelectedFile] = useState(null);
    const [base64String, setBase64String] = useState('');
    const [error, setError] = useState('');
    const [message, setMessage] = useState('');
    const userId = Number(localStorage.getItem('userId'));

    const handleFileChange = (event) => {
        const file = event.target.files[0];
        const reader = new FileReader();

        if (file) {
            reader.onloadend = () => {
                const base64data = reader.result.split(",")[1]; // Eliminar la cabecera del base64
                setBase64String(base64data);
                setSelectedFile(file);
            };
            reader.onerror = () => {
                setError("Error reading file");
            };
            reader.readAsDataURL(file);
        }
    };

    const handleUpload = async () => {
        if (!selectedFile) {
            setError('Please select a file to upload.');
            return;
        }

        if (!courseId) {
            setError('Course ID is required.');
            return;
        }

        console.log(base64String); // Verificar el contenido base64

        try {
            const response = await axios.post('http://localhost:8080/files', {
                name: selectedFile.name,
                content: base64String,
                userId: userId,
                courseId: Number(courseId),
            });

            if (response.status === 200) {
                setMessage('File uploaded successfully');
                setSelectedFile(null); // Clear the selected file after upload
                setBase64String('');
            } else {
                setError('Failed to upload file');
            }
        } catch (error) {
            if (error.response) {
                setError(`Failed to upload file: ${error.response.data.error}`);
            } else if (error.request) {
                setError('No response received from server');
            } else {
                setError(`Error in setting up request: ${error.message}`);
            }
        }
    };

    return (
        <div>
            <input type="file" accept=".txt,.pdf" onChange={handleFileChange} /> {/* Allow txt and pdf files */}
            <button onClick={handleUpload}>Upload</button>
            {error && <p style={{ color: 'red' }}>{error}</p>}
            {message && <p style={{ color: 'green' }}>{message}</p>}
        </div>
    );
};

export default FileUploadComponent;

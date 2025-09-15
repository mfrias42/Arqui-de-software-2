import React, { useState, useRef } from 'react';
import axios from 'axios';
import { useNavigate } from 'react-router-dom';
import '../assets/styles/AddCourse.css';

function AddCourse() {
    const [name, setName] = useState('');
    const [description, setDescription] = useState('');
    const [category, setCategory] = useState('');
    const [duration, setDuration] = useState('');
    const [instructorId, setInstructorId] = useState('');
    const [capacity, setCapacity] = useState('');
    const [imageBase64, setImageBase64] = useState('');
    const [imagePreview, setImagePreview] = useState('');
    const [imageName, setImageName] = useState('');
    const [error, setError] = useState('');
    const navigate = useNavigate();
    const fileInputRef = useRef();

    const handleImageChange = (e) => {
        const file = e.target.files[0];
        if (file) {
            setImageName(file.name);
            const reader = new FileReader();
            reader.onloadend = () => {
                setImageBase64(reader.result);
                setImagePreview(reader.result);
            };
            reader.readAsDataURL(file);
        }
    };

    const handleFileButtonClick = (e) => {
        e.preventDefault();
        fileInputRef.current.click();
    };

    const handleSubmit = async (e) => {
        e.preventDefault();

        if (!name || !description || !category || !duration || !instructorId || !capacity) {
            setError('Todos los campos son obligatorios');
            return;
        }

        const parsedInstructorId = parseInt(instructorId);
        const parsedCapacity = parseInt(capacity);
        if (isNaN(parsedInstructorId) || isNaN(parsedCapacity)) {
            setError('ID del instructor y capacidad deben ser números');
            return;
        }

        const courseData = {
            name,
            description,
            category,
            duration,
            instructor_id: parsedInstructorId,
            capacity: parsedCapacity,
            imageBase64: imageBase64 || undefined
        };

        const token = localStorage.getItem('token');
        try {
            await axios.post('http://localhost:8080/courses', courseData, {
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${token}`
                }
            });
            setError('');
            alert('Curso agregado con éxito');
            navigate('/manage-courses');
        } catch (error) {
            setError('Error al agregar curso: ' + (error.response?.data?.message || error.message));
        }
    };

    return (
        <div className="add-course-bg">
            <div className="add-course-container">
                <button className="back-button" onClick={() => navigate('/home')}>Volver</button>
                <h1>Agregar nuevo curso</h1>
                {error && <p className="error-message">{error}</p>}
                <form onSubmit={handleSubmit} className="add-course-form">
                    <input type="text" value={name} onChange={e => setName(e.target.value)} placeholder="Nombre del curso" />
                    <input type="text" value={description} onChange={e => setDescription(e.target.value)} placeholder="Descripción" />
                    <input type="text" value={category} onChange={e => setCategory(e.target.value)} placeholder="Categoría" />
                    <input type="text" value={duration} onChange={e => setDuration(e.target.value)} placeholder="Duración" />
                    <input 
                        type="number" 
                        value={instructorId} 
                        onChange={e => setInstructorId(e.target.value)} 
                        placeholder="ID del instructor" 
                        min="0"
                    />
                    <input 
                        type="number" 
                        value={capacity} 
                        onChange={e => setCapacity(e.target.value)} 
                        placeholder="Capacidad del curso" 
                        min="1"
                    />
                    <label style={{marginTop: '10px', fontWeight: 500}}>Imagen del curso (opcional):</label>
                    <input 
                        type="file" 
                        accept="image/*" 
                        onChange={handleImageChange} 
                        ref={fileInputRef}
                        style={{display: 'none'}}
                    />
                    <button 
                        type="button"
                        className="file-upload-btn"
                        onClick={handleFileButtonClick}
                        style={{marginBottom: 8, display: 'block'}}
                    >
                        Seleccionar imagen
                    </button>
                    {imageName && <div style={{fontSize: '0.95rem', color: '#444', marginBottom: 8}}>{imageName}</div>}
                    {imagePreview && (
                        <img src={imagePreview} alt="Preview" style={{maxWidth: 180, margin: '10px auto', borderRadius: 8}} />
                    )}
                    <button type="submit" className="submit-button">Agregar curso</button>
                </form>
            </div>
        </div>
    );
}

export default AddCourse;

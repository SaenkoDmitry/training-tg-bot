import Home from './pages/Home';
import WorkoutPage from './pages/WorkoutPage';
import StatsPage from './pages/StatsPage';
import ProgramsPage from './pages/ProgramsPage';
import MeasurementsPage from './pages/MeasurementsPage';
import LibraryPage from './pages/LibraryPage';
import {Route, Routes} from 'react-router-dom';
import MainLayout from './components/MainLayout';
import React from "react";

import {AuthProvider} from './context/AuthContext';

const App = () => (
    <AuthProvider>
            <Routes>
                    <Route path="/" element={<MainLayout><Home/></MainLayout>}/>
                    <Route path="/workout/:id" element={<MainLayout><WorkoutPage/></MainLayout>}/>
                    <Route path="/stats" element={<MainLayout><StatsPage/></MainLayout>}/>
                    <Route path="/programs" element={<MainLayout><ProgramsPage/></MainLayout>}/>
                    <Route path="/measurements" element={<MainLayout><MeasurementsPage/></MainLayout>}/>
                    <Route path="/library" element={<MainLayout><LibraryPage/></MainLayout>}/>
            </Routes>
    </AuthProvider>
);

export default App;

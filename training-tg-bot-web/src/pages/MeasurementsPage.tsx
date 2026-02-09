import React, { useEffect, useState, useRef, useCallback } from 'react';
import './MeasurementsPage.css';

const PAGE_SIZE = 15;

const MeasurementsPage: React.FC = () => {
    const [measurements, setMeasurements] = useState<Measurement[]>([]);
    const [count, setCount] = useState<number>(0);
    const [offset, setOffset] = useState(0);
    const [loading, setLoading] = useState(false);
    const [hasMore, setHasMore] = useState(true);

    const [adding, setAdding] = useState(false);
    const [newMeasurement, setNewMeasurement] = useState<Partial<ToCreateMeasurement>>({});

    const observer = useRef<IntersectionObserver | null>(null);

    const lastRowRef = useCallback(
        (node: HTMLTableRowElement | null) => {
            if (loading) return;
            if (!hasMore) return;
            if (observer.current) observer.current.disconnect();
            observer.current = new IntersectionObserver(entries => {
                if (entries[0].isIntersecting) {
                    setOffset(prev => prev + PAGE_SIZE);
                }
            });
            if (node) observer.current.observe(node);
        },
        [loading, hasMore]
    );

    useEffect(() => {
        if (!hasMore) return;

        setLoading(true);

        fetch(`/api/measurements?offset=${offset}&limit=${PAGE_SIZE}`)
            .then(res => res.json())
            .then(data => {
                setMeasurements(prev => {
                    const newIds = new Set(prev.map(m => m.id));
                    const uniqueItems = data.items.filter(m => !newIds.has(m.id));
                    return [...prev, ...uniqueItems];
                });
                setCount(data.count);
                if (offset + data.items.length >= data.count) {
                    setHasMore(false);
                }
            })
            .finally(() => setLoading(false));
    }, [offset]);

    const handleInputChange = (field: keyof ToCreateMeasurement, value: string) => {
        setNewMeasurement(prev => ({
            ...prev,
            [field]: Number(value)
        }));
    };

    const handleSaveNewMeasurement = () => {
        fetch('/api/measurements', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(newMeasurement),
        })
            .then(res => res.json())
            .then(data => {
                setMeasurements(prev => [data, ...prev]); // добавляем первым
                setCount(prev => prev + 1);
                setAdding(false);
                setNewMeasurement({});
            });
    };

    const handleCancelNewMeasurement = () => {
        setAdding(false);
        setNewMeasurement({});
    };

    return (
        <div className="measurements-page">
            <h1>Замеры</h1>

            {!adding && (
                <button className="add-button" onClick={() => setAdding(true)}>
                    Добавить новое измерение
                </button>
            )}

            <div className="table-wrapper" style={{ maxHeight: '70vh', overflowY: 'auto' }}>
                <table className="measurements-table">
                    <thead>
                    <tr>
                        <th>№</th>
                        <th>Дата</th>
                        <th>Плечи</th>
                        <th>Грудь</th>
                        <th>Л. рука</th>
                        <th>П. рука</th>
                        <th>Талия</th>
                        <th>Ягодицы</th>
                        <th>Л. бедро</th>
                        <th>П. бедро</th>
                        <th>Л. икра</th>
                        <th>П. икра</th>
                        <th>Вес</th>
                        <th>Действия</th>
                    </tr>
                    </thead>

                    <tbody>
                    {adding && (
                        <tr className="new-measurement-row">
                            <td>-</td>
                            <td>-</td>
                            {(['shoulders','chest','hand_left','hand_right','waist','buttocks','hip_left','hip_right','calf_left','calf_right','weight'] as (keyof ToCreateMeasurement)[]).map(field => (
                                <td key={field}>
                                    <input
                                        type="number"
                                        value={newMeasurement[field] ?? ''}
                                        onChange={e => handleInputChange(field, e.target.value)}
                                    />
                                </td>
                            ))}
                            <td>
                                <button onClick={handleSaveNewMeasurement}>Save</button>
                                <button onClick={handleCancelNewMeasurement}>Cancel</button>
                            </td>
                        </tr>
                    )}

                    {measurements.map((m, idx) => (
                        <tr key={m.id} ref={idx === measurements.length - 1 ? lastRowRef : null}>
                            <td>{idx + 1}</td>
                            <td>{m.created_at}</td>
                            <td>{m.shoulders}</td>
                            <td>{m.chest}</td>
                            <td>{m.hand_left}</td>
                            <td>{m.hand_right}</td>
                            <td>{m.waist}</td>
                            <td>{m.buttocks}</td>
                            <td>{m.hip_left}</td>
                            <td>{m.hip_right}</td>
                            <td>{m.calf_left}</td>
                            <td>{m.calf_right}</td>
                            <td className="weight">{m.weight}</td>
                            <td></td>
                        </tr>
                    ))}

                    {loading && (
                        <tr>
                            <td colSpan={14}>Загрузка...</td>
                        </tr>
                    )}
                    </tbody>
                </table>
            </div>
        </div>
    );
};

export default MeasurementsPage;

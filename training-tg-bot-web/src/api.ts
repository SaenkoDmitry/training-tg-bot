export async function getWorkouts(userId: number) {
    // Для теста возвращаем фиктивные данные
    return [
        { id: 1, title: 'Приседания', reps: 15 },
        { id: 2, title: 'Отжимания', reps: 20 },
        { id: 3, title: 'Подтягивания', reps: 10 }
    ];
}

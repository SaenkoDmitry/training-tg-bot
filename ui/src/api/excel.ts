import {apiBlob} from "./client.ts";

export const downloadExcelWorkouts = async () => {
    const blob = await apiBlob("/api/excel/workouts", {method: "GET"});

    const url = window.URL.createObjectURL(blob);
    const a = document.createElement("a");
    a.href = url;
    a.download = "workouts.xlsx";
    document.body.appendChild(a);
    a.click();
    a.remove();
    window.URL.revokeObjectURL(url);
};

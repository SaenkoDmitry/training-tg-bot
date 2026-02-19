import React from "react";
import {
    Bird,
    Bus,
    Cat,
    Cherry,
    ChessKing,
    ChessKnight,
    ChessQueen,
    Dog,
    Flame,
    Heart,
    Lollipop,
    Panda,
    Rabbit,
    Smile,
    Squirrel,
    Turtle
} from "lucide-react";

export const ICONS = {
    Smile,
    Panda,
    Cat,
    Dog,
    Rabbit,
    Bird,
    Squirrel,
    Lollipop,
    Flame,
    Heart,
    ChessKing,
    ChessQueen,
    ChessKnight,
    Turtle,
    Bus,
    Cherry,
};

export type IconName = keyof typeof ICONS;

interface Props {
    selected: IconName;
    onSelect: (name: IconName) => void;
    onClose: () => void;
}

const IconPicker: React.FC<Props> = ({selected, onSelect, onClose}) => {
    return (
        <div
            style={{
                position: "fixed",
                inset: 0,
                background: "rgba(0,0,0,0.4)",
                display: "flex",
                justifyContent: "center",
                alignItems: "center",
                zIndex: 1000,
            }}
            onClick={onClose}
        >
            <div
                onClick={(e) => e.stopPropagation()}
                style={{
                    background: "#fff",
                    padding: "1.5rem",
                    borderRadius: 20,
                    width: 320,
                    display: "grid",
                    gridTemplateColumns: "repeat(4, 1fr)",
                    gap: 16,
                    boxShadow: "0 10px 30px rgba(0,0,0,0.15)",
                }}
            >
                {Object.entries(ICONS).map(([name, Icon]) => (
                    <div
                        key={name}
                        onClick={() => onSelect(name as IconName)}
                        style={{
                            cursor: "pointer",
                            padding: 12,
                            borderRadius: 14,
                            display: "flex",
                            justifyContent: "center",
                            transition: "all 0.2s ease",
                            background:
                                selected === name
                                    ? "rgba(0,0,0,0.06)"
                                    : "transparent",
                            border:
                                selected === name
                                    ? "2px solid #000"
                                    : "2px solid transparent",
                        }}
                    >
                        <Icon size={28}/>
                    </div>
                ))}
            </div>
        </div>
    );
};

export default IconPicker;

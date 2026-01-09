package constants

import "fmt"

const (
	LegsAndShouldersWorkoutID   = "legs_and_shoulders"
	LegsAndShouldersWorkoutName = "🦵 Ноги & плечи"

	BackAndBicepsWorkoutID   = "back_and_biceps"
	BackAndBicepsWorkoutName = "🏋️‍♂️ Спина & бицепсы"

	ChestAndTricepsID   = "chest_and_triceps"
	ChestAndTricepsName = "🫀 Грудь & трицепсы"
)

const (
	LegsID   = "legs"
	LegsName = "🦵 Ноги"

	PressID   = "press"
	PressName = "📰 Пресс"

	DeltasID   = "deltas"
	DeltasName = "δ Дельты"

	BackID   = "back"
	BackName = "🏋 Спина"

	ChestID   = "chest"
	ChestName = "🫀 Грудь"

	BicepsID   = "biceps"
	BicepsName = "💪 Бицепс"

	TricepsID   = "triceps"
	TricepsName = "💪🏻 Трицепс"

	CardioID   = "cardio"
	CardioName = "🏃 Кардио"
)

var (
	Groups = map[string]string{
		LegsID:    LegsName,
		PressID:   PressName,
		DeltasID:  DeltasName,
		BackID:    BackName,
		BicepsID:  BicepsName,
		ChestID:   ChestName,
		TricepsID: TricepsName,
		CardioID:  CardioName,
	}
)

const (
	FrontSurfaceOfTheThigh = "передняя поверхность бедра"
	BackSurfaceOfTheThigh  = "задняя поверхность бедра"

	RectusAbdominisMuscle = "прямая мышца живота"

	BackDeltoidMuscle   = "задняя дельтовидная мышца"
	MiddleDeltoidMuscle = "средняя дельтовидная мышца"

	LatissimusDorsiMuscle = "широчайшая мышца спины"

	PectoralMuscles = "грудные мышцы"

	BicepsBrachiiMuscle         = "двуглавая мышца плеча"
	TricepsShoulderMuscle       = "трехглавая мышца плеча"
	BrachialisAndshoulderMuscle = "брахиалис + плечевая мышца"
)

type ExerciseObj struct {
	ID            int
	Name          string
	Url           string
	Type          string
	Accent        string
	RestInSeconds int
}

func (e *ExerciseObj) GetName() string {
	if e == nil {
		return ""
	}
	return e.Name
}

func (e *ExerciseObj) GetAccent() string {
	if e == nil {
		return ""
	}
	return e.Accent
}

func (e *ExerciseObj) GetHint() string {
	if e == nil {
		return ""
	}
	if e.Url == "" {
		return ""
	}
	return WrapYandexLink(e.Url)
}

func WrapYandexLink(url string) string {
	return fmt.Sprintf("\n<a href=\"%s\"><b>⚠️Техника выполнения:</b></a>", url)
}

var AllExercises = map[string]*ExerciseObj{
	// Ноги
	ExtensionOfLowerLegWhileSitting: {
		ID:            1,
		Name:          "Разгибание голени сидя",
		Url:           "https://disk.yandex.ru/i/nevoPFhHbc8l8g",
		Type:          LegsID,
		RestInSeconds: 120,
		Accent:        FrontSurfaceOfTheThigh,
	},
	FlexionOfLowerLegWhileSitting: {
		ID:            2,
		Name:          "Сгибание голени сидя",
		Url:           "https://disk.yandex.ru/i/PqkWBjSSNvH-Vg",
		Type:          LegsID,
		RestInSeconds: 120,
		Accent:        BackSurfaceOfTheThigh,
	},
	PlatformLegPress: {
		ID:            3,
		Name:          "Жим платформы ногами",
		Url:           "https://disk.yandex.ru/i/UsaW3YjvDRWm3w",
		Type:          LegsID,
		RestInSeconds: 180,
		Accent:        FrontSurfaceOfTheThigh,
	},
	// Пресс
	LiftingLegsAtTheElbow: {
		ID:            4,
		Name:          "Подъем ног в висе на локтях",
		Url:           "https://disk.yandex.ru/i/pkAxqVWTe4L_Xw",
		Type:          PressID,
		RestInSeconds: 90,
		Accent:        RectusAbdominisMuscle,
	},
	// Дельты
	ReverseDilutionsInThePectoral: {
		ID:            5,
		Name:          "Обратные разведения в пек-дек",
		Url:           "https://disk.yandex.ru/i/9lYAV1wr3VjWQQ",
		Type:          DeltasID,
		RestInSeconds: 120,
		Accent:        BackDeltoidMuscle,
	},
	ExtensionOfBarbell: {
		ID:            6,
		Name:          "Протяжка штанги",
		Url:           "https://disk.yandex.ru/i/0aaEdn5IBOI6zQ",
		Type:          DeltasID,
		RestInSeconds: 120,
		Accent:        MiddleDeltoidMuscle,
	},
	//Спина
	PullUpInTheGravitronWithAWideGrip: {
		ID:            7,
		Name:          "Подтягивание в гравитроне широким хватом",
		Url:           "https://disk.yandex.ru/i/jp52K-HTe86iLA",
		Type:          BackID,
		RestInSeconds: 120,
		Accent:        LatissimusDorsiMuscle,
	},
	VerticalTractionInALeverSimulator: {
		ID:            8,
		Name:          "Вертикальная тяга в рычажном тренажере",
		Url:           "https://disk.yandex.ru/i/x6qRCfJBGA7tEQ",
		Type:          BackID,
		RestInSeconds: 120,
		Accent:        LatissimusDorsiMuscle,
	},
	HorizontalDeadliftInABlockSimulatorWithAnEmphasisOnTheChest: {
		ID:            9,
		Name:          "Горизонтальная тяга в блочном тренажере с упором в грудь",
		Url:           "https://disk.yandex.ru/i/DnyJDcPaJLUyCg",
		Type:          BackID,
		RestInSeconds: 120,
		Accent:        LatissimusDorsiMuscle,
	},
	DumbbellDeadliftWithEmphasisOnTheBench: {
		ID:            10,
		Name:          "Тяга гантели с упором в скамью",
		Url:           "https://disk.yandex.ru/i/mU9TIaxDPV6nXw",
		Type:          BackID,
		RestInSeconds: 120,
		Accent:        LatissimusDorsiMuscle,
	},
	// Бицепс
	ArmFlexionWithDumbbellSupination: {
		ID:            11,
		Name:          "Сгибание рук с супинацией гантелями",
		Url:           "https://disk.yandex.ru/i/LWBPrSeWvxNiUw",
		Type:          BicepsID,
		RestInSeconds: 120,
		Accent:        BicepsBrachiiMuscle,
	},
	HammerBendsWithDumbbells: {
		ID:            12,
		Name:          "Молотковые сгибания с гантелями",
		Url:           "https://disk.yandex.ru/i/OvY5i3YGxyi6gw",
		Type:          BicepsID,
		RestInSeconds: 120,
		Accent:        BrachialisAndshoulderMuscle,
	},
	// Грудные
	BenchPressWithAWideGrip: {
		ID:            13,
		Name:          "Жим лежа широким хватом",
		Url:           "https://disk.yandex.ru/i/w2FIsYgqMQ-RPA",
		Type:          ChestID,
		RestInSeconds: 180,
		Accent:        PectoralMuscles,
	},
	HorizontalBenchPressInTheTechnoGymSimulator: {
		ID:            14,
		Name:          "Жим горизонтально в тренажере TechnoGym",
		Url:           "https://disk.yandex.ru/i/vyDhCyusHft5VQ",
		Type:          ChestID,
		RestInSeconds: 120,
		Accent:        PectoralMuscles,
	},
	BringingArmsTogetherInTheButterflySimulator: {
		ID:            15,
		Name:          "Сведение рук в тренажере бабочка",
		Url:           "https://disk.yandex.ru/i/JADkjm4tiUsAdQ",
		Type:          ChestID,
		RestInSeconds: 120,
		Accent:        PectoralMuscles,
	},
	// Трицепс
	FrenchBenchPressWithDumbbells: {
		ID:            16,
		Name:          "Французский жим с гантелями лежа",
		Url:           "https://disk.yandex.ru/i/9KPxatabvDYy8g",
		Type:          TricepsID,
		RestInSeconds: 120,
		Accent:        TricepsShoulderMuscle,
	},
	ExtensionOfTricepsFromTheUpperBlockWithARopeHandle: {
		ID:            17,
		Name:          "Разгибание на трицепс с верхнего блока канатной рукоятью",
		Url:           "https://disk.yandex.ru/i/sG8luvJYQWNgyg",
		Type:          TricepsID,
		RestInSeconds: 120,
		Accent:        TricepsShoulderMuscle,
	},
	// Кардио
	Walking: {
		ID:   18,
		Name: "Ходьба",
		Url:  "",
		Type: CardioID,
	},
	RunningOnTrack: {
		ID:   19,
		Name: "Бег на дорожке",
		Url:  "",
		Type: CardioID,
	},
	RunningOnMechanicalTrack: {
		ID:   20,
		Name: "Бег на механической дорожке",
		Url:  "",
		Type: CardioID,
	},
	Bicycle: {
		ID:   21,
		Name: "Велосипед",
		Url:  "",
		Type: CardioID,
	},
}

const (
	// Ноги
	ExtensionOfLowerLegWhileSitting = "extension_of_lower_leg_while_sitting"
	FlexionOfLowerLegWhileSitting   = "flexion_of_lower_leg_while_sitting"
	PlatformLegPress                = "platform_leg_press"

	// Пресс
	LiftingLegsAtTheElbow = "lifting_legs_at_the_elbow"

	// Дельты
	ReverseDilutionsInThePectoral = "reverse_dilutions_in_the_pectoral"
	ExtensionOfBarbell            = "extension_of_barbell"

	//Спина
	PullUpInTheGravitronWithAWideGrip                           = "pull_up_in_the_gravitron_with_a_wide_grip"
	VerticalTractionInALeverSimulator                           = "vertical_traction_in_a_lever_simulator"
	HorizontalDeadliftInABlockSimulatorWithAnEmphasisOnTheChest = "horizontal_deadlift_in_a_block_simulator_with_an_emphasis_on_the_chest"
	DumbbellDeadliftWithEmphasisOnTheBench                      = "dumbbell_deadlift_with_emphasis_on_the_bench"

	// Бицепс
	ArmFlexionWithDumbbellSupination = "arm_flexion_with_dumbbell_supination"
	HammerBendsWithDumbbells         = "hammer_bends_with_dumbbells"

	// Грудные
	BenchPressWithAWideGrip                     = "bench_press_with_a_wide_grip"
	HorizontalBenchPressInTheTechnoGymSimulator = "horizontal_bench_press_in_the_techno_gym_simulator"
	BringingArmsTogetherInTheButterflySimulator = "bringing_arms_together_in_the_butterfly_simulator"

	// Трицепс
	FrenchBenchPressWithDumbbells                      = "french_bench_press_with_dumbbells"
	ExtensionOfTricepsFromTheUpperBlockWithARopeHandle = "extension_of_triceps_from_the_upper_block_with_a_rope_handle"

	// Кардио
	Walking                  = "walking"
	RunningOnTrack           = "running_on_track"
	RunningOnMechanicalTrack = "running_on_mechanical_track"
	Bicycle                  = "bicycle"
)

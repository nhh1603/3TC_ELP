module Main exposing (..)

import Browser
import Html exposing (..)
import Html.Attributes exposing (style, placeholder, value)
import Html.Events exposing (..)
import Http
import Json.Decode exposing (Decoder, map2, field, string, list)
import Json.Decode exposing (Error(..))
import Set exposing (..)
import Random
import Array exposing (..)



-- MAIN


main : Program () Model Msg
main =
  Browser.element
    { init = init
    , update = update
    , subscriptions = subscriptions
    , view = view
    }



-- MODEL


type alias Model =
    { status : Status
    , stringListChosenWords : List String
    , randomWord : String
    , guess : String
    , statusGuess : StatusGuess
    , showAnswer: Bool
    }


type StatusGuess
    = None 
    | FailureGuess
    | SuccessGuess


type Status
    = Failure
    | Loading
    | SuccessListWord (List Word)


type alias Word =
    { word : String
    , meanings : List Meaning
    }


type alias Meaning =
    { partOfSpeech : String
    , definitions : List String
    }


modelInit : Model
modelInit =
    { status = Loading
    , stringListChosenWords = []
    , randomWord = ""
    , guess = ""
    , statusGuess = None
    , showAnswer = False
    }


wordInit : Word
wordInit = 
    { word = "hi"
    , meanings = 
        [ { partOfSpeech = "hi", definitions = ["hi"]}
        , { partOfSpeech = "hi", definitions = ["hi"]}
        ]
    }


init : () -> ( Model, Cmd Msg )
init _ =
    ( modelInit, fetchListChosenWords )



-- UPDATE


type Msg
  = GotStringList (Result Http.Error String)
  | GenerateRandom
  | RandomIndex Int
  | GotWord (Result Http.Error (List Word))
  | ChangeGuess String
  | TryGuess
  | ShowAnswer
  | HideAnswer


update : Msg -> Model -> (Model, Cmd Msg)
update msg model =
    case msg of
        GotStringList result ->
            case result of 
                Ok stringList ->
                    let 
                        list =
                            String.words stringList
                    in
                    ( { model | stringListChosenWords = list}, Random.generate RandomIndex (Random.int 0 (List.length list - 1)))
                Err _ ->
                        ( { model | status = Failure }, fetchListChosenWords )

        GenerateRandom ->
            ( { model | statusGuess = None, guess = "", showAnswer = False }, Random.generate RandomIndex (Random.int 0 (List.length model.stringListChosenWords - 1)))

        RandomIndex index ->
            let 
                selected =
                    Array.fromList model.stringListChosenWords
                        |> Array.get index
                        |> Maybe.withDefault ""
            in
                case selected of
                    "" -> ( model, Cmd.none )
                    _ ->
                        ( { model | randomWord = selected }, fetchWord selected )

        GotWord result ->
            case result of 
                Ok listWord ->
                    ( { model | status = SuccessListWord listWord }, Cmd.none)
                Err _ ->
                    ( { model | status = Failure }, Random.generate RandomIndex (Random.int 0 (List.length model.stringListChosenWords - 1)))

        ChangeGuess newGuess ->
            ( { model | guess = newGuess }, Cmd.none )

        TryGuess ->
            if model.guess == model.randomWord then
                ( { model | statusGuess = SuccessGuess, guess = "" }, Cmd.none )
            else
                ( { model | statusGuess = FailureGuess }, Cmd.none )

        ShowAnswer ->
            ( { model | showAnswer = True }, Cmd.none)

        HideAnswer ->
            ( { model | showAnswer = False }, Cmd.none)





-- SUBSCRIPTIONS


subscriptions : Model -> Sub Msg
subscriptions model =
  Sub.none



-- VIEW


view : Model -> Html Msg
view model =
    div []
        [ h1 [ style "text-align" "center" ] [ text "The Guessing Game 🤫" ]
        , viewWord model
        ]


viewWord : Model -> Html Msg
viewWord model =
    case model.status of
        Failure ->
            div []
                [ text "I could not load a random word for some reason. "
                ]

        Loading ->
            text "Loading..."

        SuccessListWord listWord ->
            div [ style "text-align" "middle" ]
                [ h2 [ style "text-align" "center" ] [ text "Guess what word it is?" ]
                , p [ style "text-align" "middle" ]
                    [ ul [] (List.map showWord listWord)
                    ]
                , p [ style "text-align" "center" ] 
                    [ input [ placeholder "Enter your guess here", value model.guess, onInput ChangeGuess ] []
                    , button [ onClick TryGuess ] [ text "Enter" ]
                    , button [ onClick GenerateRandom ] [ text "Random"] 
                    , button [ if model.showAnswer then onClick HideAnswer else onClick ShowAnswer ] 
                             [ if model.showAnswer then text "Hide Answer" else text "Show Answer" ]
                    ]
                , div []
                      [ if model.statusGuess == SuccessGuess then 
                            h3 [ style "text-align" "center" ] [ text ( "That's right 🥳! The correct word is: " ++ model.randomWord ++ ". Press 'Random' to play with another word!") ]
                        else if model.statusGuess == FailureGuess then 
                            h3 [ style "text-align" "center" ] [ text "Sorry, that's wrong 😢! You should try again or random another word if it's too hard to guess!" ]
                        else 
                            text "" 
                      ]
                , div []
                      [ if model.showAnswer then
                            h3 [ style "text-align" "center" ] [ text ("The answer is: " ++ model.randomWord) ]
                        else
                            text ""
                      ]
                ]


showWord : Word -> Html Msg
showWord word =
    div []
        [ li [] (List.map showMeaning word.meanings)
        ]


showMeaning : Meaning -> Html Msg
showMeaning meaning =
    div []
        [ li [] [ text meaning.partOfSpeech ]
        , ol [] (List.map showDef meaning.definitions)
        ]


showDef : String -> Html Msg
showDef def =
    div []
        [ li [] [ text def]]


-- HTTP


fetchWord : String -> Cmd Msg
fetchWord word =
    let 
        url =
            "https://api.dictionaryapi.dev/api/v2/entries/en/" ++ word
    in
        Http.get 
            { url = url
            , expect = Http.expectJson GotWord wordListDecoder
            }


fetchListChosenWords : Cmd Msg
fetchListChosenWords =
    Http.get
    { url = "../static/stringListChosenWords.txt"
    , expect = Http.expectString GotStringList
    }
    


wordListDecoder : Decoder (List Word)
wordListDecoder = 
    list wordDecoder


wordDecoder : Decoder Word
wordDecoder =
    map2 Word
        (field "word" string)
        (field "meanings" (list meaningDecoder))


meaningDecoder : Decoder Meaning
meaningDecoder = 
    map2 Meaning
        (field "partOfSpeech" string)
        (field "definitions" (list (field "definition" string)))

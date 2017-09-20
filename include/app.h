#if defined(NANOGUI_GLAD)
#if defined(NANOGUI_SHARED) && !defined(GLAD_GLAPI_EXPORT)
#define GLAD_GLAPI_EXPORT
#endif

#include <glad/glad.h>
#else
#if defined(__APPLE__)
#define GLFW_INCLUDE_GLCOREARB
#else
#define GL_GLEXT_PROTOTYPES
#endif
#endif

#include <GLFW/glfw3.h>

#include <nanogui/nanogui.h>

#include <memory>

namespace my
{

class Drawable
{
  public:
    virtual void draw() = 0;
    virtual void zoom(bool in, float x, float y) = 0;
    virtual void move(float dx, float dy) = 0;
};

class Screen
{
    bool clicked = false;
    bool just_clicked = false;
    float x;
    float y;
    nanogui::Screen *screen;
    Drawable *content;

  public:
    Screen(GLFWwindow *window)
    {
        screen = new nanogui::Screen();
        screen->initialize(window, true);
    }

    void drawContents()
    {
        if (content)
        {
            content->draw();
        }
    }

    nanogui::Screen *nanogui()
    {
        return screen;
    }

    void setContent(Drawable *content)
    {
        this->content = content;
    }

    bool cursorPosCallbackEvent(double x, double y);
    bool mouseButtonCallbackEvent(int button, int action, int modifiers);
    bool scrollCallbackEvent(double x, double y);
};

#pragma pack(push, 1)
struct Point
{
    float x;
    float y;
};
#pragma pack(pop)

class Fractal : public Drawable
{
    float zoomVal = ZOOM_START;
    int iterations = ITERATIONS_START;
    float reZ = RE_Z_START;
    float imZ = IM_Z_START;
    float shiftX = 0;
    float shiftY = 0;

    GLuint program;
    GLuint vbuffer;
    GLuint varrays;
    GLuint texture;

    Point points[4] = {
        {-1.0f, -1.0f},
        {1.0f, -1.0f},
        {-1.0f, 1.0f},
        {1.0f, 1.0f}};

  public:
    static constexpr float ZOOM_START = 0.3;
    static constexpr float ZOOM_K = 1.2;
    static constexpr float ZOOM_MIN = 0.1;
    static constexpr float ZOOM_MAX = 1e4;
    static constexpr float ITERATIONS_START = 40;
    static constexpr float RE_Z_START = 1;
    static constexpr float IM_Z_START = 0;
    
    Fractal();
    void draw() override;
    void zoom(bool in, float x, float y) override
    {
        shiftX -= (x - 400) / zoomVal;
        shiftY -= (y - 400) / zoomVal;

        if (in)
        {
            zoomVal /= ZOOM_K;
            if (zoomVal < ZOOM_MIN) {
                zoomVal = ZOOM_MIN;
            }
        }
        else
        {
            zoomVal *= ZOOM_K;
            if (zoomVal > ZOOM_MAX) {
                zoomVal = ZOOM_MAX;
            }
        }

        shiftX += (x - 400) / zoomVal;
        shiftY += (y - 400) / zoomVal;
    }

    void move(float dx, float dy) override
    {
        shiftX += dx / zoomVal;
        shiftY += dy / zoomVal;
    }

    void setIterations(int value)
    {
        iterations = value;
    }

    void setReZ(float value)
    {
        reZ = value;
    }

    void setImZ(float value)
    {
        imZ = value;
    }
};

class App
{
    nanogui::FormHelper *gui;
    nanogui::Window *nanoguiWindow;
    std::unique_ptr<Screen> screen;

    nanogui::Slider *iterations;
    nanogui::Slider *reZ;
    nanogui::Slider *imZ;

    std::unique_ptr<Fractal> fractal;

  public:
    void setUpGui(GLFWwindow *window);

    Screen *getScreen()
    {
        return screen.get();
    }
};
}
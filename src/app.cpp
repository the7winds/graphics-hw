#include "app.h"

#include "shader.hpp"

struct gimp_image
{
    unsigned int width;
    unsigned int height;
    unsigned int bytes_per_pixel; /* 2:RGB16, 3:RGB, 4:RGBA */
    unsigned char pixel_data[256 * 1 * 2 + 1];
};

extern const struct gimp_image gimp_texture;

my::Fractal::Fractal()
{
    program = loadProgram();

    glGenBuffers(1, &vbuffer);
    glGenVertexArrays(1, &varrays);
    glGenTextures(1, &texture);
}

void my::Fractal::draw()
{
    glBindBuffer(GL_ARRAY_BUFFER, vbuffer);
    glBufferData(GL_ARRAY_BUFFER, sizeof(points), points, GL_STATIC_DRAW);

    glUseProgram(program);

    glBindVertexArray(varrays);
    glVertexAttribPointer(0, 2, GL_FLOAT, GL_FALSE, sizeof(Point), (GLvoid *)0);
    glEnableVertexAttribArray(0);

    glBindTexture(GL_TEXTURE_1D, texture);
    glTexImage1D(GL_TEXTURE_1D, 0, GL_RGBA, gimp_texture.width, 0, GL_RGBA, GL_UNSIGNED_BYTE, gimp_texture.pixel_data);
    glTexParameteri(GL_TEXTURE_1D, GL_TEXTURE_MAG_FILTER, GL_LINEAR);
    glTexParameteri(GL_TEXTURE_1D, GL_TEXTURE_MIN_FILTER, GL_LINEAR);

    glUniform1i(glGetUniformLocation(program, "iterations"), iterations);
    glUniform1f(glGetUniformLocation(program, "reZ"), reZ);
    glUniform1f(glGetUniformLocation(program, "imZ"), imZ);
    glUniform1f(glGetUniformLocation(program, "zoom"), zoomVal);
    glUniform1f(glGetUniformLocation(program, "shiftX"), shiftX);
    glUniform1f(glGetUniformLocation(program, "shiftY"), shiftY);

    glDrawArrays(GL_TRIANGLE_STRIP, 0, 4);
}

bool my::Screen::cursorPosCallbackEvent(double x, double y)
{
    if (screen->cursorPosCallbackEvent(x, y))
    {
        return true;
    }

    if (clicked)
    {
        float dx = x - this->x;
        float dy = y - this->y;
        content->move(dx, dy);

        this->x = x;
        this->y = y;

        return true;
    }

    this->x = x;
    this->y = y;

    return false;
}

bool my::Screen::mouseButtonCallbackEvent(int button, int action, int modifiers)
{
    if (screen->mouseButtonCallbackEvent(button, action, modifiers))
    {
        return true;
    }

    if (button == 0)
    {
        clicked = action;
        just_clicked = true;
        return true;
    }

    return false;
}

bool my::Screen::scrollCallbackEvent(double x, double y)
{
    if (!screen->scrollCallbackEvent(x, y))
    {
        content->zoom(y > 0, this->x, this->y);
    }
    return true;
}

void my::App::setUpGui(GLFWwindow *window)
{
    // Create nanogui gui
    screen.reset(new Screen(window));
    fractal.reset(new Fractal());
    screen->setContent(fractal.get());

    gui = new nanogui::FormHelper(screen->nanogui());
    nanoguiWindow = gui->addWindow(Eigen::Vector2i(10, 10), "Settings");
    nanoguiWindow->setWidth(400);

    App *t = this;

    iterations = new nanogui::Slider(nanoguiWindow);
    iterations->setRange(std::make_pair(1, 800));
    iterations->setValue(Fractal::ITERATIONS_START);
    iterations->setFixedWidth(300);
    iterations->setCallback([t](float value) { t->fractal->setIterations(value); });

    reZ = new nanogui::Slider(nanoguiWindow);
    reZ->setRange(std::make_pair(-1, 1));
    reZ->setValue(Fractal::RE_Z_START);
    reZ->setCallback([t](float value) { t->fractal->setReZ(value); });

    imZ = new nanogui::Slider(nanoguiWindow);
    imZ->setRange(std::make_pair(-1, 1));
    imZ->setValue(Fractal::IM_Z_START);
    imZ->setCallback([t](float value) { t->fractal->setImZ(value); });

    gui->addWidget("iterations", iterations);
    gui->addWidget("re seed", reZ);
    gui->addWidget("im seed", imZ);

    screen->nanogui()->setVisible(true);
    screen->nanogui()->performLayout();

    nanoguiWindow->center();
}
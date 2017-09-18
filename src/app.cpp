#include "app.h"

#include "shader.hpp"

my::Fractal::Fractal()
{
    program = loadProgram();

    glGenBuffers(1, &vbuffer);
    glGenVertexArrays(1, &varrays);
}

void my::Fractal::draw()
{
    glBindBuffer(GL_ARRAY_BUFFER, vbuffer);
    glBufferData(GL_ARRAY_BUFFER, sizeof(points), points, GL_STATIC_DRAW);

    glUseProgram(program);

    glBindVertexArray(varrays);
    glVertexAttribPointer(0, 2, GL_FLOAT, GL_FALSE, sizeof(Point), (GLvoid *)0);
    glEnableVertexAttribArray(0);

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
        content->zoom(y > 0);
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
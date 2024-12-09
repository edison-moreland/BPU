`timescale 1 ns / 1 ns

`include "alu.sv"

`define assert(test_name, signal, value) \
        if (signal !== value) begin \
            $display("ASSERTION FAILED for test %s %m: %b/%d != %b/%d", test_name, signal, signal, value, value); \
            $finish; \
        end

module alu_tb #(
    parameter N = 8
);

    logic [N-1:0] a, b, c;
    logic [2:0] op;
    logic carry_in, carry_out;
    logic a_gt_b, a_eq_b, c_eq_zero;

    alu #(.N(N))
    test(
        .a(a), .b(b), .op(op), .carry_in(carry_in),
        .c(c), .carry_out(carry_out), .a_gt_b(a_gt_b),
        .a_eq_b(a_eq_b), .c_eq_zero(c_eq_zero)
    );

    initial begin

        // Addition, no carry in
        op = 3'b000;
        a = 4; b = 8;
        carry_in = 0;
        #10
        `assert("addition_no_carry", c, (a+b))
        `assert("addition_no_carry", carry_out, 0)
        `assert("addition_no_carry", a_gt_b, 0)
        `assert("addition_no_carry", a_eq_b, 0)
        `assert("addition_no_carry", c_eq_zero, 0)

        // Addition, with carry in
        op = 3'b000;
        a = 4; b = 8;
        carry_in = 1;
        #10
        `assert("addition_with_carry", c, (a+b+1))
        `assert("addition_with_carry", carry_out, 0)
        `assert("addition_with_carry", a_gt_b, 0)
        `assert("addition_with_carry", a_eq_b, 0)
        `assert("addition_with_carry", c_eq_zero, 0)

        // Addition, no carry in, carry out
        op = 3'b000;
        a = 254; b = 1;
        carry_in = 1;
        #10
        `assert("addition_carry_out", c, 0)
        `assert("addition_carry_out", carry_out, 1)
        `assert("addition_carry_out", a_gt_b, 1)
        `assert("addition_carry_out", a_eq_b, 0)
        `assert("addition_carry_out", c_eq_zero, 1)

        // Shift Right, no carry in, no carry out
        op = 3'b001;
        a = 8'b01111110; b = 0;
        carry_in = 0;
        #10
        `assert("shift_right_no_carry_in_no_carry_out", c, 8'b00111111)
        `assert("shift_right_no_carry_in_no_carry_out", carry_out, 0)
        `assert("shift_right_no_carry_in_no_carry_out", a_gt_b, 1)
        `assert("shift_right_no_carry_in_no_carry_out", a_eq_b, 0)
        `assert("shift_right_no_carry_in_no_carry_out", c_eq_zero, 0)

        // Shift Right, carry in, no carry out
        op = 3'b001;
        a = 8'b01111110; b = 0;
        carry_in = 1;
        #10
        `assert("shift_right_carry_in_no_carry_out", c, 8'b10111111)
        `assert("shift_right_carry_in_no_carry_out", carry_out, 0)
        `assert("shift_right_carry_in_no_carry_out", a_gt_b, 1)
        `assert("shift_right_carry_in_no_carry_out", a_eq_b, 0)
        `assert("shift_right_carry_in_no_carry_out", c_eq_zero, 0)

        // Shift Right, no carry in, carry out
        op = 3'b001;
        a = 8'b01111111; b = 0;
        carry_in = 0;
        #10
        `assert("shift_right_no_carry_in_carry_out", c, 8'b00111111)
        `assert("shift_right_no_carry_in_carry_out", carry_out, 1)
        `assert("shift_right_no_carry_in_carry_out", a_gt_b, 1)
        `assert("shift_right_no_carry_in_carry_out", a_eq_b, 0)
        `assert("shift_right_no_carry_in_carry_out", c_eq_zero, 0)

        // Shift Right, carry in, carry out
        op = 3'b001;
        a = 8'b01111111; b = 0;
        carry_in = 1;
        #10
        `assert("shift_right_carry_in_carry_out", c, 8'b10111111)
        `assert("shift_right_carry_in_carry_out", carry_out, 1)
        `assert("shift_right_carry_in_carry_out", a_gt_b, 1)
        `assert("shift_right_carry_in_carry_out", a_eq_b, 0)
        `assert("shift_right_carry_in_carry_out", c_eq_zero, 0)

        // Shift left, no carry in, no carry out
        op = 3'b010;
        a = 8'b01111110; b = 0;
        carry_in = 0;
        #10
        `assert("shift_left_no_carry_in_no_carry_out", c, 8'b11111100)
        `assert("shift_left_no_carry_in_no_carry_out", carry_out, 0)
        `assert("shift_left_no_carry_in_no_carry_out", a_gt_b, 1)
        `assert("shift_left_no_carry_in_no_carry_out", a_eq_b, 0)
        `assert("shift_left_no_carry_in_no_carry_out", c_eq_zero, 0)

        // Shift left, carry in, no carry out
        op = 3'b010;
        a = 8'b01111110; b = 0;
        carry_in = 1;
        #10
        `assert("shift_left_carry_in_no_carry_out", c, 8'b11111101)
        `assert("shift_left_carry_in_no_carry_out", carry_out, 0)
        `assert("shift_left_carry_in_no_carry_out", a_gt_b, 1)
        `assert("shift_left_carry_in_no_carry_out", a_eq_b, 0)
        `assert("shift_left_carry_in_no_carry_out", c_eq_zero, 0)

        // Shift left, no carry in, carry out
        op = 3'b010;
        a = 8'b11111110; b = 0;
        carry_in = 0;
        #10
        `assert("shift_left_no_carry_in_carry_out", c, 8'b11111100)
        `assert("shift_left_no_carry_in_carry_out", carry_out, 1)
        `assert("shift_left_no_carry_in_carry_out", a_gt_b, 1)
        `assert("shift_left_no_carry_in_carry_out", a_eq_b, 0)
        `assert("shift_left_no_carry_in_carry_out", c_eq_zero, 0)

        // Shift left, carry in, carry out
        op = 3'b010;
        a = 8'b11111110; b = 0;
        carry_in = 1;
        #10
        `assert("shift_left_carry_in_carry_out", c, 8'b11111101)
        `assert("shift_left_carry_in_carry_out", carry_out, 1)
        `assert("shift_left_carry_in_carry_out", a_gt_b, 1)
        `assert("shift_left_carry_in_carry_out", a_eq_b, 0)
        `assert("shift_left_carry_in_carry_out", c_eq_zero, 0)

        // Carry in won't be used from here on out
        carry_in = 0;

        // NOT
        op = 3'b011;
        a = 8'b01010101; b = 0;
        #10
        `assert("not", c, 8'b10101010)
        `assert("not", carry_out, 0)
        `assert("not", a_gt_b, 1)
        `assert("not", a_eq_b, 0)
        `assert("not", c_eq_zero, 0)

        // AND
        op = 3'b100;
        a = 8'b01010101; b = 8'b11111111;
        #10
        `assert("and", c, 8'b01010101)
        `assert("and", carry_out, 0)
        `assert("and", a_gt_b, 0)
        `assert("and", a_eq_b, 0)
        `assert("and", c_eq_zero, 0)

        // OR
        op = 3'b101;
        a = 8'b01010101; b = 8'b10101111;
        #10
        `assert("or", c, 8'b11111111)
        `assert("or", carry_out, 0)
        `assert("or", a_gt_b, 0)
        `assert("or", a_eq_b, 0)
        `assert("or", c_eq_zero, 0)

        // XOR
        op = 3'b110;
        a = 8'b01010101; b = 8'b10101111;
        #10
        `assert("xor", c, 8'b11111010)
        `assert("xor", carry_out, 0)
        `assert("xor", a_gt_b, 0)
        `assert("xor", a_eq_b, 0)
        `assert("xor", c_eq_zero, 0)

        // Compare greater
        op = 3'b111;
        a = 8; b = 4;
        #10
        `assert("compare_greater", c, 0)
        `assert("compare_greater", carry_out, 0)
        `assert("compare_greater", a_gt_b, 1)
        `assert("compare_greater", a_eq_b, 0)
        `assert("compare_greater", c_eq_zero, 1)

        // Compare not greater
        op = 3'b111;
        a = 4; b = 8;
        #10
        `assert("compare_not_greater", c, 0)
        `assert("compare_not_greater", carry_out, 0)
        `assert("compare_not_greater", a_gt_b, 0)
        `assert("compare_not_greater", a_eq_b, 0)
        `assert("compare_not_greater", c_eq_zero, 1)

        // Compare equal
        op = 3'b111;
        a = 4; b = 4;
        #10
        `assert("compare_equal", c, 0)
        `assert("compare_equal", carry_out, 0)
        `assert("compare_equal", a_gt_b, 0)
        `assert("compare_equal", a_eq_b, 1)
        `assert("compare_equal", c_eq_zero, 1)
    end

endmodule
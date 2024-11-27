`include "full_adder.sv"

//  TODO: Most of this get's optimized away. We need a visualization that works for this

// Carry select half of alu adder
module alu_adder_cs_half (
   input  logic [3:0] a, b,
   input  logic       carry_in,
   input  logic       enable,
   output logic [3:0] result,
   output logic       carry_out,
);

    logic [3:0] low_result;
    logic       low_carry;

    full_adder #(.N(4)) low_adder(
        .a(a), .b(b),
        .carry_in(0),
        .result(low_result),
        .carry_out(low_carry)
    );

    logic [3:0] high_result;
    logic       high_carry;

    full_adder #(.N(4)) high_adder(
        .a(a), .b(b),
        .carry_in(1),
        .result(high_result),
        .carry_out(high_carry)
    );

    always_comb begin
        if (enable)
            if (carry_in)
                result <= high_result;
            else
                result <= low_result;
        else
            result <= 0;

        if (carry_in)
            carry_out <= high_carry;
        else
            carry_out <= low_carry;
    end

endmodule